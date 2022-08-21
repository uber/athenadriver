// Copyright (c) 2022 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package athenadriver

import (
	"context"
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go/service/athena/athenaiface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/athena"
)

// Rows defines rows in AWS Athena ResultSet.
type Rows struct {
	athena          athenaiface.AthenaAPI
	ctx             context.Context
	queryID         string
	reachedLastPage bool
	ResultOutput    *athena.GetQueryResultsOutput
	config          *Config
	tracer          *DriverTracer
	pageCount       int64
}

// NewNonOpsRows is to create a new Rows.
func NewNonOpsRows(ctx context.Context, athenaAPI athenaiface.AthenaAPI, queryID string, driverConfig *Config,
	obs *DriverTracer) (*Rows, error) {
	r := Rows{
		athena:    athenaAPI,
		ctx:       ctx,
		queryID:   queryID,
		config:    driverConfig,
		tracer:    obs,
		pageCount: -1,
	}
	return &r, nil
}

// NewRows is to create a new Rows.
func NewRows(ctx context.Context, athenaAPI athenaiface.AthenaAPI, queryID string, driverConfig *Config,
	obs *DriverTracer) (*Rows, error) {
	r := Rows{
		athena:    athenaAPI,
		ctx:       ctx,
		queryID:   queryID,
		config:    driverConfig,
		tracer:    obs,
		pageCount: -1,
	}
	if err := r.fetchNextPage(nil); err != nil {
		return nil, err
	}
	return &r, nil
}

// Columns return Columns metadata.
func (r *Rows) Columns() []string {
	var columns []string
	for _, colInfo := range r.ResultOutput.ResultSet.ResultSetMetadata.ColumnInfo {
		columns = append(columns, *colInfo.Name)
	}
	return columns
}

// ColumnTypeDatabaseTypeName will be called by sql framework.
func (r *Rows) ColumnTypeDatabaseTypeName(index int) string {
	colInfo := r.ResultOutput.ResultSet.ResultSetMetadata.ColumnInfo[index]
	if colInfo.Type != nil {
		return *colInfo.Type
	}
	r.tracer.Scope().Counter(DriverName + ".failure.columntypedatabasetypename").Inc(1)
	r.tracer.Log(ErrorLevel, "ColumnTypeDatabaseTypeName failed", zap.Int("index", index))
	return ""
}

// Next is to get next result set page.
func (r *Rows) Next(dest []driver.Value) error {
	if r.reachedLastPage {
		return io.EOF
	}
	if len(r.ResultOutput.ResultSet.Rows) == 0 {
		if r.ResultOutput.NextToken == nil || *r.ResultOutput.NextToken == "" {
			// this means we reach the last page - no token and no rows
			r.reachedLastPage = true
			return io.EOF
		}

		if err := r.fetchNextPage(r.ResultOutput.NextToken); err != nil {
			return err
		}
		if r.reachedLastPage {
			return io.EOF
		}
	}

	// Shift to next row
	cur := r.ResultOutput.ResultSet.Rows[0]
	columns := r.ResultOutput.ResultSet.ResultSetMetadata.ColumnInfo
	if err := r.convertRow(columns, cur.Data, dest, r.config); err != nil {
		return err
	}
	r.ResultOutput.ResultSet.Rows = r.ResultOutput.ResultSet.Rows[1:]
	return nil
}

// fetchNextPage is to get next result set page with a specific token.
func (r *Rows) fetchNextPage(token *string) error {
	var err error
	r.ResultOutput, err = r.athena.GetQueryResultsWithContext(r.ctx,
		&athena.GetQueryResultsInput{
			QueryExecutionId: aws.String(r.queryID),
			NextToken:        token,
		})
	if err != nil {
		r.tracer.Scope().Counter(DriverName + ".failure.fetchnextpage.getqueryresults").Inc(1)
		r.tracer.Log(ErrorLevel, "GetQueryResults failed", zap.String("error", err.Error()))
		r.reachedLastPage = true
		return err
	}

	r.pageCount++
	// First row of the first page contains header if the query is not DDL.
	// These are also available in *athenaAPI.Row.ResultSetMetadata.
	// Sometimes Athena go API will return row data without corresponding ColumnInfo. To circumvent this situation,
	// we choose to name the column as `column` + 0-index-based number
	// One example is:
	//   input:
	//      MSCK REPAIR TABLE sampledb.elb_logs
	//   output:
	//     _col0
	//     Partitions not in metastore:    elb_logs:2015/01/01     elb_logs:2015/01/02     elb_logs:2015/01/03
	//       elb_logs:2015/01/04     elb_logs:2015/01/05     elb_logs:2015/01/06     elb_logs:2015/01/07
	if r.ResultOutput != nil &&
		r.ResultOutput.ResultSet.ResultSetMetadata != nil &&
		r.ResultOutput.ResultSet.ResultSetMetadata.ColumnInfo != nil {
		rowLen := len(r.ResultOutput.ResultSet.Rows)
		colLen := len(r.ResultOutput.ResultSet.ResultSetMetadata.ColumnInfo)
		if rowLen > 0 {
			rowColLen := len(r.ResultOutput.ResultSet.Rows[0].Data)
			if colLen < rowColLen {
				for i := 0; i < rowColLen-colLen; i++ {
					colName := "_col" + strconv.Itoa(i+colLen)
					colType := "string"
					colInfo := newColumnInfo(colName, colType)
					r.ResultOutput.ResultSet.ResultSetMetadata.ColumnInfo = append(r.ResultOutput.ResultSet.ResultSetMetadata.ColumnInfo,
						colInfo)
				}
			} else if colLen > rowColLen && rowColLen == 1 {
				for k := 0; k < rowLen; k++ {
					items := strings.Split(*r.ResultOutput.ResultSet.Rows[k].Data[0].VarCharValue, "\t")
					if len(items) == colLen {
						for i, v := range items {
							items[i] = strings.TrimSpace(v)
						}
						r.ResultOutput.ResultSet.Rows[k] = newRow(colLen, items)
					}
				}
			}
		} else if rowLen == 0 && colLen == 1 && r.ResultOutput.UpdateCount != nil {
			if *r.ResultOutput.UpdateCount > 0 {
				if *r.ResultOutput.ResultSet.ResultSetMetadata.ColumnInfo[0].Name == "rows" {
					// For DML's INSERT INTO, DDL's CTAS
					updateCount := strconv.FormatInt(*r.ResultOutput.UpdateCount, 10)
					rData := athena.Datum{VarCharValue: &updateCount}
					aRow := athena.Row{Data: []*athena.Datum{&rData}}
					r.ResultOutput.ResultSet.Rows = append(r.ResultOutput.ResultSet.Rows, &aRow)
				}
			}
		}
	}
	var rowOffset = 0
	if r.pageCount == 0 {
		rs := r.ResultOutput.ResultSet
		ci := r.ResultOutput.ResultSet.ResultSetMetadata.ColumnInfo
		i := 0
		if len(ci) > 0 && len(rs.Rows) > 0 && len(rs.Rows[0].Data) > 0 && len(rs.Rows[0].Data) == len(ci) {
			for ; i < len(ci); i++ {
				if rs.Rows[0].Data[i] == nil || rs.Rows[0].Data[i].VarCharValue == nil {
					break
				}
				if *ci[i].Name != *rs.Rows[0].Data[i].VarCharValue {
					break
				}
			}
			if i == len(ci) {
				rowOffset = 1
			}
		}
	}

	// if there is no new row, we should not continue, and this also filters out cases that Rows is nil
	if len(r.ResultOutput.ResultSet.Rows) <= rowOffset {
		r.reachedLastPage = true
		return nil
	}

	r.ResultOutput.ResultSet.Rows = r.ResultOutput.ResultSet.Rows[rowOffset:]
	return nil
}

// Close is to close Rows after reading all data.
func (r *Rows) Close() error {
	if r.ResultOutput != nil && r.ResultOutput.NextToken != nil {
		r.tracer.Log(WarnLevel, "rows close prematurely, queryID: "+r.queryID)
		r.ResultOutput = nil
	}
	r.reachedLastPage = true
	return nil
}

// convertRow is to convert data from Athena type to Golang SQL type and put them into an array of driver.Value.
func (r *Rows) convertRow(columns []*athena.ColumnInfo, rdata []*athena.Datum, ret []driver.Value,
	driverConfig *Config) error {
	for i, val := range rdata {
		if val == nil {
			return ErrAthenaNilDatum
		}
		value, err := r.athenaTypeToGoType(columns[i], val.VarCharValue, driverConfig)
		if err != nil {
			r.tracer.Log(ErrorLevel, "convertrow failed", zap.String("error", err.Error()))
			r.tracer.Scope().Counter(DriverName + ".failure.convertrow").Inc(1)
			return err
		}
		/*r.tracer.Log(DebugLevel, "TM",
			zap.String("athenaType", *columns[i].Type),
			zap.String("goType", reflect.TypeOf(value).String()),
			zap.String("str", *val.VarCharValue),
		)*/
		ret[i] = value
	}
	return nil
}

// athenaTypeToGoType converts Athena type to Golang SQL type.
// https://docs.aws.amazon.com/en_pv/athena/latest/ug/data-types.html
// https://docs.aws.amazon.com/athena/latest/ug/geospatial-input-data-formats-supported-geometry-types.html#geometry-data-types
// varbinary is undocumented above, but appears in geo query like:
//
//	SELECT ST_POINT(-74.006801, 40.705220).
//
// json is also undocumented above, but appears here https://docs.aws.amazon.com/athena/latest/ug/querying-JSON.html
// The full list is here: https://prestodb.io/docs/0.172/language/types.html
// Include ipaddress for forward compatibility.
func (r *Rows) athenaTypeToGoType(columnInfo *athena.ColumnInfo, rawValue *string, driverConfig *Config) (interface{}, error) {
	if maskedValue, masked := driverConfig.CheckColumnMasked(*columnInfo.Name); masked { // "comma ok" idiom
		return maskedValue, nil
	}
	if rawValue == nil {
		r.tracer.Scope().Counter(DriverName + ".missingvalue").Inc(1)
		r.tracer.Log(ErrorLevel, "missing data",
			zap.String("columnInfo.Name", *columnInfo.Name),
			zap.String("queryID", r.queryID),
			zap.String("workgroup", driverConfig.GetWorkgroup().Name))
		if driverConfig.IsMissingAsEmptyString() {
			return "", nil
		} else if driverConfig.IsMissingAsDefault() {
			return r.getDefaultValueForColumnType(*columnInfo.Type), nil
		}
		r.tracer.Scope().Counter(DriverName + ".failure.convertvalue.config").Inc(1)
		r.tracer.Log(ErrorLevel, "missing data", zap.String("columnInfo.Name", *columnInfo.Name))
		return nil, fmt.Errorf("Missing data at column " + *columnInfo.Name)
	}
	val := *rawValue
	// https://stackoverflow.com/questions/30299649/parse-string-to-specific-type-of-int-int8-int16-int32-int64
	// https://prestodb.io/docs/current/language/types.html#integer
	var err error
	var i int64
	var f float64
	switch *columnInfo.Type {
	case "tinyint":
		// strconv.ParseInt() behavior is to return (int64(0), err)
		// which is not as good as just return (nil, err)
		if i, err = strconv.ParseInt(val, 10, 8); err != nil {
			return nil, err
		}
		return int8(i), nil
	case "smallint":
		if i, err = strconv.ParseInt(val, 10, 16); err != nil {
			return nil, err
		}
		return int16(i), nil
	case "integer":
		if i, err = strconv.ParseInt(val, 10, 32); err != nil {
			return nil, err
		}
		return int32(i), nil
	case "bigint":
		if i, err = strconv.ParseInt(val, 10, 64); err != nil {
			return nil, err
		}
		return i, nil
	case "float", "real":
		if f, err = strconv.ParseFloat(val, 32); err != nil {
			return nil, err
		}
		return float32(f), nil
	case "double":
		if f, err = strconv.ParseFloat(val, 64); err != nil {
			return nil, err
		}
		return f, nil
	// for binary, we assume all chars are 0 or 1; for json,
	// we assume the json syntax is correct. Leave to caller to verify it.
	case "json", "char", "varchar", "varbinary", "row", "string", "binary",
		"struct", "interval year to month", "interval day to second", "decimal",
		"ipaddress", "array", "map", "unknown":
		return val, nil
	case "boolean":
		if val == "true" {
			return true, nil
		} else if val == "false" {
			return false, nil
		}
		r.tracer.Scope().Counter(DriverName + ".failure.convertvalue.boolean").Inc(1)
		r.tracer.Log(ErrorLevel, "boolean data error", zap.String("val", val))
		return nil, fmt.Errorf("unknown value `%s` for boolean", val)
	case "date", "time", "time with time zone", "timestamp", "timestamp with time zone":
		vv, err := scanTime(val)
		if !vv.Valid {
			r.tracer.Scope().Counter(DriverName + ".failure.convertvalue." +
				"time").Inc(1)
			r.tracer.Log(ErrorLevel, "time data error",
				zap.String("val", val),
				zap.String("type", *columnInfo.Type))
			return nil, err
		}
		return vv.Time, err
	default:
		r.tracer.Scope().Counter(DriverName + ".failure.convertvalue.type").Inc(1)
		r.tracer.Log(ErrorLevel, "column data type error", zap.String("columnInfo.Type", *columnInfo.Type))
		return nil, fmt.Errorf("unknown type `%s` with value %s", *columnInfo.Type, val)
	}
}

// getDefaultValueForColumnType is used internally by athenaTypeToGoType to get default value for a column type.
// This is helpful when column has missing value and we want to display it anyway.
func (r *Rows) getDefaultValueForColumnType(athenaType string) interface{} {
	switch athenaType {
	case "tinyint", "smallint", "integer", "bigint":
		return 0
	case "boolean":
		return false
	case "float", "double", "real":
		return 0.0
	case "date", "time", "time with time zone", "timestamp", "timestamp with time zone":
		return time.Time{}
	case "json", "char", "varchar", "varbinary", "row", "string", "binary",
		"struct", "interval year to month", "interval day to second", "decimal",
		"ipaddress", "array", "map", "unknown":
		return ""
	default:
		r.tracer.Scope().Counter(DriverName + ".failure.defaultvalueforcolumntype.type").Inc(1)
		r.tracer.Log(ErrorLevel, "column data type error", zap.String("columnInfo.Type", athenaType))
		return ""
	}
}
