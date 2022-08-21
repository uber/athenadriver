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
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/xwb1989/sqlparser"
)

// OutputStyles are all the styles we can choose to print query result
var OutputStyles = [...]string{"StyleDefault", "StyleBold", "StyleColoredBright", "StyleColoredDark",
	"StyleColoredBlackOnBlueWhite", "StyleColoredBlackOnCyanWhite", "StyleColoredBlackOnGreenWhite",
	"StyleColoredBlackOnMagentaWhite", "StyleColoredBlackOnYellowWhite", "StyleColoredBlackOnRedWhite",
	"StyleColoredBlueWhiteOnBlack", "StyleColoredCyanWhiteOnBlack", "StyleColoredGreenWhiteOnBlack",
	"StyleColoredMagentaWhiteOnBlack", "StyleColoredRedWhiteOnBlack", "StyleColoredYellowWhiteOnBlack",
	"StyleDouble", "StyleLight", "StyleRounded",
}

// OutputFormats are all the formats we can choose to print query result
var OutputFormats = [...]string{"csv", "html", "markdown", "table"}

func scanNullString(v interface{}) (sql.NullString, error) {
	if v == nil {
		return sql.NullString{}, nil
	}
	vv, ok := v.(string)
	if !ok {
		return sql.NullString{},
			fmt.Errorf("cannot convert %v (%T) to string", v, v)
	}
	return sql.NullString{Valid: true, String: vv}, nil
}

func mockRowsToSQLRows(mockRows *sqlmock.Rows) *sql.Rows {
	db, mock, _ := sqlmock.New()
	mock.ExpectQuery("SELECT_OK").WillReturnRows(mockRows)
	rows, _ := db.Query("SELECT_OK")
	return rows
}

// ColsToCSV is a convenient function to convert columns of sql.Rows to CSV format.
func ColsToCSV(rows *sql.Rows) string {
	if rows == nil {
		return ""
	}
	columns, _ := rows.Columns()
	s := ""
	for i, v := range columns {
		s += v
		if i != len(columns)-1 {
			s += ","
		} else {
			s += "\n"
		}
	}
	return s
}

// RowsToCSV is to convert rows of sql.Rows to CSV format.
func RowsToCSV(rows *sql.Rows) string {
	if rows == nil {
		return ""
	}
	columns, _ := rows.Columns()
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	records := make([][]string, 0)
	for rows.Next() {
		rawResult := make([][]byte, len(columns))
		row := make([]interface{}, len(columns))
		for i := range rawResult {
			row[i] = &rawResult[i] // pointers to each string in the interface slice
		}
		// We don't consider malformed rows
		_ = rows.Scan(row...)
		s := make([]string, len(columns))
		for i, cell := range rawResult {
			s[i] = string(cell)
		}
		records = append(records, s)
	}
	csvWriter.WriteAll(records)
	return buf.String()
}

// ColsRowsToCSV is a convenient function to convert columns and rows of sql.Rows to CSV format.
func ColsRowsToCSV(rows *sql.Rows) string {
	s := ColsToCSV(rows)
	r := RowsToCSV(rows)
	return s + r
}

func getTableStyle(style string) table.Style {
	switch style {
	case "StyleColoredBright":
		return table.StyleColoredBright
	case "StyleBold":
		return table.StyleBold
	case "StyleColoredDark":
		return table.StyleColoredDark
	case "StyleColoredBlackOnBlueWhite":
		return table.StyleColoredBlackOnBlueWhite
	case "StyleColoredBlackOnCyanWhite":
		return table.StyleColoredBlackOnCyanWhite
	case "StyleColoredBlackOnGreenWhite":
		return table.StyleColoredBlackOnGreenWhite
	case "StyleColoredBlackOnMagentaWhite":
		return table.StyleColoredBlackOnMagentaWhite
	case "StyleColoredBlackOnYellowWhite":
		return table.StyleColoredBlackOnYellowWhite
	case "StyleColoredBlackOnRedWhite":
		return table.StyleColoredBlackOnRedWhite
	case "StyleColoredBlueWhiteOnBlack":
		return table.StyleColoredBlueWhiteOnBlack
	case "StyleColoredCyanWhiteOnBlack":
		return table.StyleColoredCyanWhiteOnBlack
	case "StyleColoredGreenWhiteOnBlack":
		return table.StyleColoredGreenWhiteOnBlack
	case "StyleColoredMagentaWhiteOnBlack":
		return table.StyleColoredMagentaWhiteOnBlack
	case "StyleColoredRedWhiteOnBlack":
		return table.StyleColoredRedWhiteOnBlack
	case "StyleColoredYellowWhiteOnBlack":
		return table.StyleColoredYellowWhiteOnBlack
	case "StyleDouble":
		return table.StyleDouble
	case "StyleLight":
		return table.StyleLight
	case "StyleRounded":
		return table.StyleRounded
	}
	return table.StyleDefault
}

func renderTable(renderType string, w table.Writer) string {
	switch renderType {
	case "markdown":
		return w.RenderMarkdown()
	case "table":
		return w.Render()
	case "html":
		return w.RenderHTML()
	}
	return w.RenderCSV()
}

// PrettyPrintSQLRows is to print rows beautifully
func PrettyPrintSQLRows(rows *sql.Rows, style string, render string, page int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if rows == nil {
		return
	}
	columns, _ := rows.Columns()
	for rows.Next() {
		rawResult := make([][]byte, len(columns))
		row := make([]interface{}, len(columns))
		for i := range rawResult {
			row[i] = &rawResult[i] // pointers to each string in the interface slice
		}
		// We don't consider malformed rows
		_ = rows.Scan(row...)
		s := make(table.Row, len(columns))
		for i, cell := range rawResult {
			s[i] = string(cell)
		}
		t.AppendRow(s)
	}
	t.SetPageSize(page)
	t.SetStyle(getTableStyle(style))
	renderTable(render, t)
}

// PrettyPrintSQLColsRows is to print rows beautifully with header
func PrettyPrintSQLColsRows(rows *sql.Rows, style string, render string, page int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if rows == nil {
		return
	}
	columns, _ := rows.Columns()
	if columns != nil && len(columns) > 0 {
		myrow := make(table.Row, len(columns))
		for i, c := range columns {
			myrow[i] = c
		}
		t.AppendHeader(myrow)
	}
	for rows.Next() {
		rawResult := make([][]byte, len(columns))
		row := make([]interface{}, len(columns))
		for i := range rawResult {
			row[i] = &rawResult[i] // pointers to each string in the interface slice
		}
		// We don't consider malformed rows
		_ = rows.Scan(row...)
		s := make(table.Row, len(columns))
		for i, cell := range rawResult {
			s[i] = string(cell)
		}
		t.AppendRow(s)
	}
	t.SetPageSize(page)
	t.SetStyle(getTableStyle(style))
	renderTable(render, t)
}

// PrettyPrintCSV is to print rows in CSV format with default style
func PrettyPrintCSV(rows *sql.Rows) {
	PrettyPrintSQLColsRows(rows, "StyleDefault", "csv", 1024)
}

// PrettyPrintMD is to print rows in markdown format with default style
func PrettyPrintMD(rows *sql.Rows) {
	PrettyPrintSQLColsRows(rows, "StyleDefault", "markdown", 1024)
}

// PrettyPrintFancy is to print rows in table format with fancy style
func PrettyPrintFancy(rows *sql.Rows) {
	PrettyPrintSQLColsRows(rows, "StyleColoredGreenWhiteOnBlack", "table", 1024)
}

// colInFirstPage is to check if this is a SELECT or VALUES statement.
// Some Sample Queries are like:
//
// USING FUNCTION predict_customer_registration(age INTEGER)
//
//	RETURNS DOUBLE TYPE
//	SAGEMAKER_INVOKE_ENDPOINT WITH (sagemaker_endpoint = 'xgboost-2019-09-20-04-49-29-303')
//
// SELECT predict_customer_registration(age) AS probability_of_enrolling, customer_id
//
//	FROM "sampledb"."ml_test_dataset"
//	WHERE predict_customer_registration(age) < 0.5;
//
// USING FUNCTION decompress(col1 VARCHAR)
//
//	RETURNS VARCHAR TYPE
//	LAMBDA_INVOKE WITH (lambda_name = 'MyAthenaUDFLambda')
//
// SELECT
//
//	decompress('ewLLinKzEsPyXdKdc7PLShKLS5OTQEAUrEH9w==');
//
// WITH
// dataset AS (
//
//	SELECT
//	  ARRAY ['hello', 'amazon', 'athena'] AS words,
//	  ARRAY ['hi', 'alexa'] AS alexa
//
// )
// SELECT concat(words, alexa) AS welcome_msg FROM dataset
func colInFirstPage(query string) bool {
	nQuery := strings.TrimSpace(strings.ToLower(query))
	return strings.Index(nQuery, "select") == 0 ||
		strings.Index(nQuery, "using") == 0 ||
		strings.Index(nQuery, "with") == 0 ||
		strings.Index(nQuery, "values") == 0
}

func isReadOnlyStatement(query string) bool {
	nQuery := strings.TrimSpace(strings.ToLower(query))
	return strings.Index(nQuery, "select") == 0 ||
		strings.Index(nQuery, "using") == 0 ||
		strings.Index(nQuery, "with") == 0 ||
		strings.Index(nQuery, "desc") == 0 ||
		strings.Index(nQuery, "show") == 0 ||
		IsQID(query)
}

func isInsertStatement(query string) bool {
	nQuery := strings.TrimSpace(strings.ToLower(query))
	return strings.Index(nQuery, "insert") == 0
}

func newColumnInfo(colName string, colType interface{}) *athena.ColumnInfo {
	caseSensitive := false
	catalogName := "hive"
	nullable := "UNKNOWN"
	precision := int64(19)
	scale := int64(0)
	schemaName := ""
	tableName := ""
	if colType == nil {
		return &athena.ColumnInfo{
			CaseSensitive: &caseSensitive,
			CatalogName:   &catalogName,
			Label:         &colName,
			Name:          &colName,
			Nullable:      &nullable,
			Precision:     &precision,
			Scale:         &scale,
			SchemaName:    &schemaName,
			TableName:     &tableName,
			Type:          nil,
		}
	}
	ct := colType.(string)
	return &athena.ColumnInfo{
		CaseSensitive: &caseSensitive,
		CatalogName:   &catalogName,
		Label:         &colName,
		Name:          &colName,
		Nullable:      &nullable,
		Precision:     &precision,
		Scale:         &scale,
		SchemaName:    &schemaName,
		TableName:     &tableName,
		Type:          &ct,
	}
}

func newRow(colLen int, rData []string) *athena.Row {
	var nData = make([]*athena.Datum, colLen)
	for i := 0; i < colLen; i++ {
		nData[i] = &athena.Datum{VarCharValue: &rData[i]}
	}
	return &athena.Row{
		Data: nData,
	}
}

func randString(l int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	s := make([]byte, l)
	for i := 0; i < l; i++ {
		s[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(s)
}

func randomInt64(min int64, max int64) int64 {
	return min + rand.Int63n(max-min)
}

// https://golang.org/ref/spec#Numeric_types
func randInt8() *string {
	s := strconv.Itoa(int(randomInt64(math.MinInt8, math.MaxInt8)))
	return &s
}

func randInt16() *string {
	s := strconv.Itoa(int(randomInt64(math.MinInt16, math.MaxInt16)))
	return &s
}

func randInt() *string {
	s := strconv.Itoa(int(randomInt64(math.MinInt32, math.MaxInt32)))
	return &s
}

func randUInt64() *string {
	s := strconv.FormatUint(rand.Uint64(), 10)
	return &s
}

func randFloat32() *string {
	s := strconv.FormatFloat(rand.Float64(), 'f', 6, 32)
	return &s
}

func randFloat64() *string {
	s := strconv.FormatFloat(rand.Float64(), 'f', 6, 64)
	return &s
}

func randStr() *string {
	s := randString(rand.Intn(10))
	return &s
}

func randBool() *string {
	if rand.Intn(10)%2 == 0 {
		s := "true"
		return &s
	}
	s := "false"
	return &s
}

func randDate() *string {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	s := time.Unix(sec, 0).Format(DateUniXFormat)
	return &s
}

func randTimeStamp() *string {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	s := time.Unix(sec, 0).Format(TimestampUniXFormat)
	return &s
}

func genHeaderRow(columns []*athena.ColumnInfo) *athena.Row {
	colLen := len(columns)
	rData := make([]string, colLen)
	for i := 0; i < colLen; i++ {
		rData[i] = *columns[i].Name
	}
	return newRow(colLen, rData)
}

// randRow generates a row with random data aligned with type information in
// athena.ColumnInfo
func randRow(columns []*athena.ColumnInfo) *athena.Row {
	colLen := len(columns)
	row := &athena.Row{
		Data: make([]*athena.Datum, colLen),
	}
	for j := 0; j < colLen; j++ {
		if columns[j].Type == nil {
			s := "a\tb"
			row.Data[j] = &athena.Datum{VarCharValue: &s}
			continue
		}
		switch *columns[j].Type {
		case "tinyint":
			row.Data[j] = &athena.Datum{VarCharValue: randInt8()}
		case "smallint":
			row.Data[j] = &athena.Datum{VarCharValue: randInt16()}
		case "integer":
			row.Data[j] = &athena.Datum{VarCharValue: randInt()}
		case "bigint":
			row.Data[j] = &athena.Datum{VarCharValue: randUInt64()}
		case "float", "real":
			row.Data[j] = &athena.Datum{VarCharValue: randFloat32()}
		case "double":
			row.Data[j] = &athena.Datum{VarCharValue: randFloat64()}
		case "json", "char", "varchar", "varbinary", "row", "string", "binary",
			"struct", "interval year to month", "interval day to second", "decimal",
			"ipaddress", "array", "map", "unknown":
			row.Data[j] = &athena.Datum{VarCharValue: randStr()}
		case "boolean":
			row.Data[j] = &athena.Datum{VarCharValue: randBool()}
		case "date":
			row.Data[j] = &athena.Datum{VarCharValue: randDate()}
		case "time", "time with time zone", "timestamp with time zone":
			row.Data[j] = &athena.Datum{VarCharValue: randTimeStamp()}
		case "timestamp":
			row.Data[j] = &athena.Datum{VarCharValue: randTimeStamp()}
		default:
			row.Data[j] = &athena.Datum{VarCharValue: randStr()}
		}
	}
	return row
}

func missingDataRow(columns []*athena.ColumnInfo) *athena.Row {
	colLen := len(columns)
	row := &athena.Row{
		Data: make([]*athena.Datum, colLen),
	}
	for j := 0; j < colLen; j++ {
		switch *columns[j].Type {
		case "integer":
			row.Data[j] = &athena.Datum{VarCharValue: nil}
		default:
			row.Data[j] = nil
		}
	}
	return row
}

func genRow(rowData []*string) *athena.Row {
	row := &athena.Row{
		Data: make([]*athena.Datum, len(rowData)),
	}
	for i := 0; i < len(rowData); i++ {
		row.Data[i] = &athena.Datum{VarCharValue: rowData[i]}
	}
	return row
}

// columnTypes must be from one of AthenaColumnTypes
func newHeaderResultPage(columnNames []*string, columnTypes []string, rowsData [][]*string) *athena.GetQueryResultsOutput {
	columns := make([]*athena.ColumnInfo, len(columnNames))
	for i := 0; i < len(columnNames); i++ {
		columns[i] = newColumnInfo(*columnNames[i], columnTypes[i])
	}
	rowLen := len(rowsData)
	rows := make([]*athena.Row, rowLen+1)
	rows[0] = genHeaderRow(columns)
	for i := 1; i < rowLen+1; i++ {
		rows[i] = genRow(rowsData[i-1])
	}
	return &athena.GetQueryResultsOutput{
		NextToken: nil,
		ResultSet: &athena.ResultSet{
			ResultSetMetadata: &athena.ResultSetMetadata{
				ColumnInfo: columns,
			},
			Rows: rows,
		},
	}
}

func newHeaderlessResultPage(columnNames []*string, columnTypes []string, rowsData [][]*string) *athena.GetQueryResultsOutput {
	columns := make([]*athena.ColumnInfo, len(columnNames))
	for i := 0; i < len(columnNames); i++ {
		columns[i] = newColumnInfo(*columnNames[i], columnTypes[i])
	}
	rowLen := len(rowsData)
	rows := make([]*athena.Row, rowLen)
	for i := 0; i < rowLen; i++ {
		rows[i] = genRow(rowsData[i])
	}
	return &athena.GetQueryResultsOutput{
		NextToken: nil,
		ResultSet: &athena.ResultSet{
			ResultSetMetadata: &athena.ResultSetMetadata{
				ColumnInfo: columns,
			},
			Rows: rows,
		},
	}
}

func newRandomHeaderResultPage(columns []*athena.ColumnInfo, nextToken *string,
	rowLen int) *athena.GetQueryResultsOutput {
	rows := make([]*athena.Row, rowLen)
	rows[0] = genHeaderRow(columns)
	for i := 1; i < rowLen; i++ {
		rows[i] = randRow(columns)
	}
	return &athena.GetQueryResultsOutput{
		NextToken: nextToken,
		ResultSet: &athena.ResultSet{
			ResultSetMetadata: &athena.ResultSetMetadata{
				ColumnInfo: columns,
			},
			Rows: rows,
		},
	}
}

func newRandomHeaderlessResultPage(columns []*athena.ColumnInfo, nextToken *string,
	rowLen int) *athena.GetQueryResultsOutput {
	rows := make([]*athena.Row, rowLen)
	for i := 0; i < rowLen; i++ {
		rows[i] = randRow(columns)
	}
	return &athena.GetQueryResultsOutput{
		NextToken: nextToken,
		ResultSet: &athena.ResultSet{
			ResultSetMetadata: &athena.ResultSetMetadata{
				ColumnInfo: columns,
			},
			Rows: rows,
		},
	}
}

// escapeBytesBackslash escapes []byte with backslashes (\)
// This escapes the contents of a string (provided as []byte) by adding backslashes before special
// characters, and turning others into specific escape sequences, such as
// turning newlines into \n and null bytes into \0.
//
// \xNN notation to define a string constant holding some peculiar byte values.
// (Of course, bytes range from hexadecimal values 00 through FF, inclusive.)
func escapeBytesBackslash(buf, v []byte) []byte {
	pos := len(buf)
	buf = reserveBuffer(buf, len(v)*2)

	for _, c := range v {
		switch c {
		case '\x00':
			buf[pos] = '\\'
			buf[pos+1] = '0'
			pos += 2
		case '\n':
			buf[pos] = '\\'
			buf[pos+1] = 'n'
			pos += 2
		case '\r':
			buf[pos] = '\\'
			buf[pos+1] = 'r'
			pos += 2
		case '\x1a':
			buf[pos] = '\\'
			buf[pos+1] = 'Z'
			pos += 2
		case '\'':
			buf[pos] = '\\'
			buf[pos+1] = '\''
			pos += 2
		case '"':
			buf[pos] = '\\'
			buf[pos+1] = '"'
			pos += 2
		case '\\':
			buf[pos] = '\\'
			buf[pos+1] = '\\'
			pos += 2
		default:
			buf[pos] = c
			pos++
		}
	}

	return buf[:pos]
}

// escapeStringBackslash is similar to escapeBytesBackslash but for string.
func escapeStringBackslash(buf []byte, v string) []byte {
	return escapeBytesBackslash(buf, []byte(v))
}

// reserveBuffer checks cap(buf) and expand buffer to len(buf) + appendSize.
// If cap(buf) is not enough, reallocate new buffer.
func reserveBuffer(buf []byte, appendSize int) []byte {
	newSize := len(buf) + appendSize
	if cap(buf) < newSize {
		// Grow buffer exponentially
		newBuf := make([]byte, len(buf)*2+appendSize)
		copy(newBuf, buf)
		buf = newBuf
	}
	return buf[:newSize]
}

func namedValueToValue(named []driver.NamedValue) []driver.Value {
	args := make([]driver.Value, len(named))
	for n, param := range named {
		args[n] = param.Value
	}
	return args
}

func valueToNamedValue(args []driver.Value) []driver.NamedValue {
	nameValues := make([]driver.NamedValue, len(args))
	for i := 0; i < len(args); i++ {
		nameValues[i].Value = args[i]
		nameValues[i].Ordinal = i + 1
	}
	return nameValues
}

func isQueryTimeOut(startOfStartQueryExecution time.Time, queryType string, serviceLimitOverride *ServiceLimitOverride) bool {
	ddlQueryTimeout := DDLQueryTimeout
	dmlQueryTimeout := DMLQueryTimeout
	if serviceLimitOverride != nil {
		if serviceLimitOverride.GetDDLQueryTimeout() > 0 {
			ddlQueryTimeout = serviceLimitOverride.GetDDLQueryTimeout()
		}
		if serviceLimitOverride.GetDMLQueryTimeout() > 0 {
			dmlQueryTimeout = serviceLimitOverride.GetDMLQueryTimeout()
		}
	}
	switch queryType {
	case "DDL":
		return time.Since(startOfStartQueryExecution) >
			time.Duration(ddlQueryTimeout)*time.Second
	case "DML":
		return time.Since(startOfStartQueryExecution) >
			time.Duration(dmlQueryTimeout)*time.Second
	case "UTILITY":
		return time.Since(startOfStartQueryExecution) >
			time.Duration(dmlQueryTimeout)*time.Second
	case "TIMEOUT_NOW":
		return true
	default:
		return time.Since(startOfStartQueryExecution) >
			time.Duration(ddlQueryTimeout)*time.Second
	}
}

// isQueryValid is to check the validity of Query, now only string length check.
// https://docs.aws.amazon.com/athena/latest/ug/service-limits.html
func isQueryValid(query string) bool {
	return len(query) < MAXQueryStringLength && len(query) > 4
}

// GetFromEnvVal is to get environmental variable value by keys.
// The return value is from whichever key is set according to the order in the slice.
func GetFromEnvVal(keys []string) string {
	for _, k := range keys {
		if v := os.Getenv(k); len(v) != 0 {
			return v
		}
	}
	return ""
}

// printCost is to print query cost
// https://aws.amazon.com/athena/pricing/
// getCost of 10MB: 5 / (1024. * 1024.) * 10 = 4.76837158203125e-05
func printCost(o *athena.GetQueryExecutionOutput) {
	if o == nil || o.QueryExecution == nil || o.QueryExecution.Statistics == nil {
		println("query cost: 0.0 USD, scanned data: 0 B, qid: NA")
		return
	}
	dataScannedBytes := o.QueryExecution.Statistics.DataScannedInBytes
	if dataScannedBytes == nil {
		println("query cost: 0.0 USD, scanned data: 0 B, qid: NA")
	} else if *dataScannedBytes == 0 {
		println("query cost: 0.0 USD, scanned data: 0 B, qid: " + *o.QueryExecution.QueryExecutionId)
	} else if *dataScannedBytes < 10*1024*1024 {
		fmt.Printf("query cost: %.20f USD, scanned data: %d B, qid: %s\n",
			getCost(*dataScannedBytes),
			*dataScannedBytes,
			*o.QueryExecution.QueryExecutionId)
	} else {
		fmt.Printf("query cost: %.20f USD, scanned data: %d B, qid: %s\n",
			getCost(*dataScannedBytes),
			*dataScannedBytes,
			*o.QueryExecution.QueryExecutionId)
	}
}

// getCost is return the USD cost upon data scanned in Bytes
// https://aws.amazon.com/athena/pricing/
func getCost(data int64) float64 {
	if data == 0 {
		return 0.0
	} else if data < int64(10*1024*1024) {
		return getPrice10MB()
	} else {
		return float64(data) * getPriceOneByte()
	}
}

var multiLineCommentPattern = regexp.MustCompile(`\/\*(.*)\*/\s*`)
var oneLineCommentPattern = regexp.MustCompile(`(^\-\-[^\n]+|\s--[^\n]+)`)
var getTableNamePattern = regexp.MustCompile(`(?i)\s+(?:from|join)\s+([\w.]+)`)
var dualPattern = regexp.MustCompile(`from dual`)
var qIDPattern = regexp.MustCompile(`^[0-9a-f-]{36}$`)

// GetTableNamesInQuery is a pessimistic function to return tables involved in query in format of DB.TABLE
// https://regoio.herokuapp.com/
// https://golang.org/pkg/regexp/syntax/
func GetTableNamesInQuery(query string) map[string]bool {
	query = multiLineCommentPattern.ReplaceAllString(query, "")
	query = oneLineCommentPattern.ReplaceAllString(query, "")
	matchedResults := getTableNamePattern.FindAllStringSubmatch(query, -1)
	tables := map[string]bool{}
	for _, matchedTableName := range matchedResults {
		if len(matchedTableName) == 2 {
			if strings.IndexByte(matchedTableName[1], '.') == -1 {
				tables[DefaultDBName+"."+matchedTableName[1]] = true
			} else {
				tables[matchedTableName[1]] = true
			}
		}
	}
	return tables
}

// GetTidySQL is to return a tidy SQL string
func GetTidySQL(query string) string {
	query = multiLineCommentPattern.ReplaceAllString(query, "")
	query = oneLineCommentPattern.ReplaceAllString(query, "")
	stmt, err := sqlparser.Parse(query)
	if err == nil {
		q := sqlparser.String(stmt)
		// OtherRead represents a DESCRIBE, or EXPLAIN statement.
		// OtherAdmin represents a misc statement that relies on ADMIN privileges.
		if q == "otherread" || q == "otheradmin" || strings.Contains(q, " '$path' ") {
			return strings.Trim(query, " ")
		}
		query = dualPattern.ReplaceAllString(q, "")
	}
	return strings.Trim(query, " ")
}

// IsQID is to check if a query string is a Query ID
// the hexadecimal Athena query ID like a44f8e61-4cbb-429a-b7ab-bea2c4a5caed
// https://aws.amazon.com/premiumsupport/knowledge-center/access-download-athena-query-results/
func IsQID(q string) bool {
	return qIDPattern.MatchString(q)
}
