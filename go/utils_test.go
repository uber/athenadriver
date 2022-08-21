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
	"database/sql"
	"database/sql/driver"
	"math"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/stretchr/testify/assert"
)

func TestScanNullString(t *testing.T) {
	s, e := scanNullString(nil)
	assert.Nil(t, e)
	assert.Equal(t, s, sql.NullString{})

	s, e = scanNullString("nil")
	assert.Nil(t, e)
	assert.True(t, s.Valid)

	s, e = scanNullString(1)
	assert.NotNil(t, e)
	assert.False(t, s.Valid)
}

func TestColsToCSV(t *testing.T) {
	sqlRows := sqlmock.NewRows([]string{"one", "two", "three"})
	rows := mockRowsToSQLRows(sqlRows)
	expected := ColsToCSV(rows)
	assert.Equal(t, expected, "one,two,three\n")
	assert.Equal(t, "", ColsToCSV(nil))
}

func TestRowsToCSV(t *testing.T) {
	sqlRows := sqlmock.NewRows([]string{"one", "two", "three"})
	sqlRows.AddRow("1", "2", "3")
	rows := mockRowsToSQLRows(sqlRows)
	expected := RowsToCSV(rows)
	assert.Equal(t, expected, "1,2,3\n")

	s := RowsToCSV(nil)
	assert.Equal(t, "", s)
}

func TestColsRowsToCSV(t *testing.T) {
	sqlRows := sqlmock.NewRows([]string{"one", "two", "three"})
	sqlRows.AddRow("1", "2", "3")
	rows := mockRowsToSQLRows(sqlRows)
	expected := ColsRowsToCSV(rows)
	assert.Equal(t, expected, "one,two,three\n1,2,3\n")
}

func TestPrettyPrintSQLRows(t *testing.T) {
	sqlRows := sqlmock.NewRows([]string{"one", "two", "three"})
	sqlRows.AddRow("1", "2", "3")
	sqlRows.AddRow("a", "b", "c")
	sqlRows.AddRow("hello", "world", "athenadriver")
	rows := mockRowsToSQLRows(sqlRows)
	for _, s := range OutputStyles {
		for _, r := range OutputFormats {
			PrettyPrintSQLRows(rows, s, r, 2)
		}
	}
	PrettyPrintMD(rows)
	PrettyPrintCSV(rows)
	PrettyPrintFancy(rows)
}

func TestPrettyPrintSQLColsRows(t *testing.T) {
	sqlRows := sqlmock.NewRows([]string{"one", "two", "three"})
	sqlRows.AddRow("1", "2", "3")
	sqlRows.AddRow("a", "b", "c")
	sqlRows.AddRow("hello", "world", "athenadriver")
	rows := mockRowsToSQLRows(sqlRows)
	PrettyPrintSQLColsRows(rows, "StyleColoredBlackOnBlueWhite", "default", 2)
	PrettyPrintSQLRows(nil, "StyleColoredBlackOnBlueWhite", "default", 2)
}

func TestPrettyPrintSQLColsRows2(t *testing.T) {
	sqlRows := sqlmock.NewRows([]string{"one", "two", "three"})
	sqlRows.AddRow("1", "2", "3")
	sqlRows.AddRow("a", "b", "c")
	sqlRows.AddRow("hello", "world", "athenadriver")
	rows := mockRowsToSQLRows(sqlRows)
	PrettyPrintSQLColsRows(rows, "StyleColoredBlackOnCyanWhite", "default", 0)
	PrettyPrintSQLColsRows(nil, "StyleColoredBlackOnCyanWhite", "default", 0)
}

func TestIsSelectStatement(t *testing.T) {
	assert.True(t, colInFirstPage("SELECT"))
	assert.True(t, colInFirstPage(" SELECT"))
	assert.True(t, colInFirstPage("select"))
}

func TestIsInsetStatement(t *testing.T) {
	assert.True(t, isInsertStatement("INSERT"))
	assert.True(t, isInsertStatement("     INSERT"))
	assert.True(t, isInsertStatement("insert"))
}

func TestRandInt8(t *testing.T) {
	s := randInt8()
	i, err := strconv.ParseInt(*s, 10, 8)
	assert.True(t, math.MinInt8 <= i && i <= math.MaxInt8)
	assert.Nil(t, err)
}

func TestRandInt16(t *testing.T) {
	s := randInt16()
	i, err := strconv.ParseInt(*s, 10, 16)
	assert.True(t, math.MinInt16 <= i && i <= math.MaxInt16)
	assert.Nil(t, err)
}

func TestRandInt(t *testing.T) {
	s := randInt()
	i, err := strconv.ParseInt(*s, 10, 32)
	assert.True(t, math.MinInt32 <= i && i <= math.MaxInt32)
	assert.Nil(t, err)
}

func TestRandInt64(t *testing.T) {
	s := randUInt64()
	_, err := strconv.ParseUint(*s, 10, 64)
	assert.Nil(t, err)
}

func TestRandFloat32(t *testing.T) {
	s := randFloat32()
	i, err := strconv.ParseFloat(*s, 32)
	assert.True(t, math.SmallestNonzeroFloat32 <= i && i <= math.MaxFloat32)
	assert.Nil(t, err)
}

func TestRandFloat64(t *testing.T) {
	s := randFloat64()
	i, err := strconv.ParseFloat(*s, 64)
	assert.True(t, math.SmallestNonzeroFloat64 <= i && i <= math.MaxFloat64)
	assert.Nil(t, err)
}

func TestRandRow(t *testing.T) {
	c1 := newColumnInfo("c1", nil)
	r := randRow([]*athena.ColumnInfo{c1})
	assert.Equal(t, len(r.Data), 1)
	assert.Equal(t, *r.Data[0].VarCharValue, "a\tb")

	for _, ty := range AthenaColumnTypes {
		c1 := newColumnInfo("c1", ty)
		r := randRow([]*athena.ColumnInfo{c1})
		assert.Equal(t, len(r.Data), 1)
	}
}

func TestNamedValueToValue(t *testing.T) {
	dn := driver.NamedValue{
		Name: "abc",
	}
	d := []driver.NamedValue{
		dn,
	}
	v := namedValueToValue(d)
	assert.Equal(t, len(v), 1)
}

type aType struct {
	S string
}

func TestValueToNamedValue(t *testing.T) {
	dn := aType{
		S: "abc",
	}
	d := []driver.Value{
		dn,
	}
	v := valueToNamedValue(d)
	assert.Equal(t, len(v), 1)
	assert.True(t, v[0].Name == "")
	assert.True(t, v[0].Ordinal == 1)
	assert.True(t, v[0].Value.(aType).S == "abc")
}

func TestIsQueryTimeOut(t *testing.T) {
	assert.False(t, isQueryTimeOut(time.Now(), athena.StatementTypeDdl, nil))
	assert.False(t, isQueryTimeOut(time.Now(), athena.StatementTypeDml, nil))
	assert.False(t, isQueryTimeOut(time.Now(), athena.StatementTypeUtility, nil))
	now := time.Now()
	OneHourAgo := now.Add(-3600 * time.Second)
	assert.True(t, isQueryTimeOut(OneHourAgo, athena.StatementTypeDml, nil))
	assert.False(t, isQueryTimeOut(OneHourAgo, athena.StatementTypeDdl, nil))
	assert.False(t, isQueryTimeOut(OneHourAgo, "UNKNOWN", nil))

	testConf := NewServiceLimitOverride()
	testConf.SetDMLQueryTimeout(65 * 60) // 65 minutes
	assert.False(t, isQueryTimeOut(OneHourAgo, athena.StatementTypeDml, testConf))

	testConf.SetDDLQueryTimeout(30 * 60) // 30 minutes
	assert.True(t, isQueryTimeOut(OneHourAgo, athena.StatementTypeDdl, testConf))
	assert.True(t, isQueryTimeOut(OneHourAgo, "UNKNOWN", testConf))
}

func TestEscapeBytesBackslash(t *testing.T) {
	r := escapeBytesBackslash([]byte{}, []byte{'\x00'})
	assert.Equal(t, string(r), "\\0")

	r = escapeBytesBackslash([]byte{}, []byte{'\n'})
	assert.Equal(t, string(r), "\\n")

	r = escapeBytesBackslash([]byte{}, []byte{'\r'})
	assert.Equal(t, string(r), "\\r")

	r = escapeBytesBackslash([]byte{}, []byte{'\x1a'})
	assert.Equal(t, string(r), "\\Z")

	r = escapeBytesBackslash([]byte{}, []byte{'\''})
	assert.Equal(t, string(r), `\'`)

	r = escapeBytesBackslash([]byte{}, []byte{'"'})
	assert.Equal(t, string(r), `\"`)

	r = escapeBytesBackslash([]byte{}, []byte{'\\'})
	assert.Equal(t, string(r), `\\`)

	r = escapeBytesBackslash([]byte{}, []byte{'x'})
	assert.Equal(t, string(r), `x`)
}

func TestGetFromEnvVal(t *testing.T) {
	os.Setenv("henrywu_test", "1")
	assert.Equal(t, GetFromEnvVal([]string{"henrywu_test"}), "1")
	assert.Equal(t, GetFromEnvVal([]string{"wufuheng", "henrywu_test"}), "1")
	os.Unsetenv("henrywu_test")
	assert.Equal(t, GetFromEnvVal([]string{"henrywu_test"}), "")
}

func TestPrintCost(t *testing.T) {
	ping := "SELECTExecContext_OK_QID"
	stat := athena.QueryExecutionStateSucceeded
	o := &athena.GetQueryExecutionOutput{
		QueryExecution: &athena.QueryExecution{
			Query:            &ping,
			QueryExecutionId: &ping,
			Status: &athena.QueryExecutionStatus{
				State: &stat,
			},
			Statistics: &athena.QueryExecutionStatistics{
				DataScannedInBytes: nil,
			},
		},
	}
	printCost(nil)
	printCost(&athena.GetQueryExecutionOutput{
		QueryExecution: nil,
	})
	printCost(&athena.GetQueryExecutionOutput{
		QueryExecution: &athena.QueryExecution{
			Query:            &ping,
			QueryExecutionId: &ping,
			Status: &athena.QueryExecutionStatus{
				State: &stat,
			},
			Statistics: nil,
		},
	})
	printCost(o)
	cost := int64(123)
	o.QueryExecution.Statistics.DataScannedInBytes = &cost
	printCost(o)
	cost = int64(12345678123456)
	o.QueryExecution.Statistics.DataScannedInBytes = &cost
	printCost(o)
	cost = int64(0)
	o.QueryExecution.Statistics.DataScannedInBytes = &cost
	printCost(o)
}

func TestUilts_GetTableNamesInQuery(t *testing.T) {
	query := "SELECT * from abc"
	tableNames := GetTableNamesInQuery(query)
	assert.Len(t, tableNames, 1)

	query = "SELECT 1"
	tableNames = GetTableNamesInQuery(query)
	assert.Len(t, tableNames, 0)

	query = "CREATE TABLE testme3 WITH (format = 'TEXTFILE', " +
		"external_location = 's3://external-location-henrywu/testme3_2', " +
		"partitioned_by = ARRAY['ssl_protocol']) AS SELECT * FROM sampledb.elb_logs"
	tableNames = GetTableNamesInQuery(query)
	assert.Len(t, tableNames, 1)

	query = "SELECT * from sampledb.abc"
	tableNames = GetTableNamesInQuery(query)
	assert.Len(t, tableNames, 1)

	query = "WITH employee AS (SELECT * FROM Employees)\nSELECT * FROM employee WHERE ID < 20\nUNION ALL\nSELECT * FROM employee WHERE Sex = 'M'"
	tableNames = GetTableNamesInQuery(query)
	assert.Len(t, tableNames, 2)

	query = "SELECT Orders.OrderID, Customers.CustomerName, Orders.OrderDate\nFROM Orders\n" +
		"INNER JOIN Customers ON Orders.CustomerID=Customers.CustomerID where a='from abc'"
	tableNames = GetTableNamesInQuery(query)
	assert.Len(t, tableNames, 2)

	// This is a case which can fail our function, but we are fine with it since we are pessimistic
	query = "SELECT Orders.OrderID, Customers.CustomerName, Orders.OrderDate\nFROM Orders\n" +
		"INNER JOIN Customers ON Orders.CustomerID=Customers.CustomerID where a='select * from abc where d=1'"
	tableNames = GetTableNamesInQuery(query)
	assert.Len(t, tableNames, 3)

	query = "select\n  sum(cost) as cost,\n  account_id,\n  case\n    user_mmowner" +
		"\n    when '' then 'Unknown'\n    else user_mmowner\n  end as user_mmowner," +
		"\n  case\n    user_owner\n    when '' then 'Unknown'\n    else user_owner" +
		"\n  end as user_owner,\n  case\n    user_team\n    when '' then 'Unknown'" +
		"\n    else user_team\n  end as user_team,\n  case\n    user_organization" +
		"\n    when '' then 'Unknown'\n    else user_organization\n  end as user_organization," +
		"\n  usage_date,\n  date_format(usage_date, '%Y-%m') as month_year\nfrom\nmytable\nwhere" +
		"\n  date_format(usage_date, '%Y-%m') >= \"2019-01\"\n  and product_name = 'AmazonAthena'\ngroup by" +
		"\n  account_id,\n  user_mmowner,\n  user_owner,\n  user_team," +
		"\n  user_organization,\n  usage_date"
	tableNames = GetTableNamesInQuery(query)
	assert.Len(t, tableNames, 1)

	query = "/* QuickSight 12345678-1234-5678-812c-4a02cbae96b8 */\nSELECT \"key\", \"bbbb\"\n" +
		"FROM (SELECT \n*,\ncount(*) as aaaa\nFROM (\nSELECT \n key\n" +
		"FROM adfads.asdfsdf.adfasdfasd\nWHERE regexp_like(assda, ',.*presto') = false \n" +
		"        AND requesturi_operation = 'GET' \n        AND bucket = 'atg-rlogs' \n" +
		"        AND key like 'manifests%' \n        AND httpstatus = '200' \n" +
		"        AND CAST( date_format(date_parse(rpad(requestdatetime, 11, '-'), '%d/%b/%Y'), '%Y-%m-%d') AS DATE )" +
		" >= (current_date - interval '7' day) \n        AND split_part(aaaa, '/', 2) not " +
		"like 'xxx'\n) subquery\nGROUP BY key ORDER BY ddddd asc LIMIT 10) " +
		"AS \"asdadsfasdfa\""
	tableNames = GetTableNamesInQuery(query)
	assert.Len(t, tableNames, 1)

	query = "-- What percentage of each node does each non-daemon pod use, and who owns it, if anyone from henry?\n" +
		"SELECT by_resource.cluster_id, node_object.name AS node_name, instance_type, pricing_model, region, zone, " +
		"namespace, by_resource.name AS pod_name, avg(resource_percent) AS responsibility_pct, owner\nFROM (\n" +
		"  SELECT pods.cluster_id, pods.uid, object.name, namespace, pods.node_uid, instance_type," +
		" req.name AS resource_name, req.quantity,\n    100 * req.quantity / sum(req.quantity) OVER w AS resource_percent\n" +
		"  FROM (\n    SELECT cluster_id, uid, node_uid\n    FROM pod_active\n    WHERE as_of <= '2019-12-02 10:00:00'\n" +
		"    UNION ALL -- We know these sets are disjoint.\n    SELECT cluster_id, uid, node_uid\n" +
		"    FROM pod_active_during\n    WHERE during @> '2019-12-02 10:00:00'::timestamp\n  ) pods\n" +
		"  NATURAL JOIN object\n  NATURAL JOIN namespaced_object\n  JOIN pod_resource_request req USING (cluster_id, uid)\n" +
		"  JOIN node_machine\n    ON node_machine.cluster_id = pods.cluster_id\n    AND node_machine.uid = pods.node_uid\n" +
		"  WHERE NOT EXISTS\n    (SELECT 1\n     FROM (\n       SELECT owner_uid AS uid\n       FROM owning_controller\n" +
		"       WHERE as_of <= '2019-12-02 10:00:00'\n         AND cluster_id = pods.cluster_id\n         AND uid = pods.uid\n" +
		"       UNION ALL -- We know these sets are disjoint.\n       SELECT owner_uid AS uid\n" +
		"       FROM owning_controller_during\n       WHERE during @> '2019-12-02 10:00:00'::timestamp\n" +
		"         AND cluster_id = pods.cluster_id\n         AND uid = pods.uid\n       LIMIT 1\n     ) owner\n" +
		"     JOIN object_type_meta\n       ON cluster_id = pods.cluster_id\n       AND object_type_meta.uid = owner.uid\n" +
		"       WHERE api_group = 'apps'\n         AND version = 'v1'\n         AND kind = 'DaemonSet'\n  )\n" +
		"  WINDOW w AS (PARTITION BY pods.cluster_id, node_uid, req.name)\n) by_resource\nJOIN object node_object\n" +
		"  ON node_object.cluster_id = by_resource.cluster_id\n  AND node_object.uid = by_resource.node_uid\n" +
		"LEFT OUTER JOIN node_pricing\n  ON node_pricing.cluster_id = by_resource.cluster_id\n" +
		"  AND node_pricing.uid = by_resource.node_uid\nLEFT OUTER JOIN node_topology\n" +
		"  ON node_topology.cluster_id = by_resource.cluster_id\n  AND node_topology.uid = by_resource.node_uid\n" +
		"LEFT OUTER JOIN (\n  SELECT cluster_id, uid, owner\n  FROM object_owner\n  WHERE as_of <= '2019-12-02 10:00:00'\n" +
		"  UNION ALL -- We know these sets are disjoint.\n  SELECT cluster_id, uid, owner\n  FROM object_owner_during\n" +
		"  WHERE during @> '2019-12-02 10:00:00'::timestamp\n) pod_owners\n  ON pod_owners.cluster_id = by_resource.cluster_id\n" +
		"  AND pod_owners.uid = by_resource.uid\nGROUP BY by_resource.cluster_id, node_name, instance_type, pricing_model, " +
		"region, zone, namespace, pod_name, owner;"
	tableNames = GetTableNamesInQuery(query)
	assert.Len(t, tableNames, 13)

	query = "/*What percentage of each node does each non-daemon pod use, and who owns it, if anyone from henry?*/\n" +
		"SELECT by_resource.cluster_id, node_object.name AS node_name, instance_type, pricing_model, region, zone, " +
		"namespace, by_resource.name AS pod_name, avg(resource_percent) AS responsibility_pct, owner\nFROM (\n" +
		"  SELECT pods.cluster_id, pods.uid, object.name, namespace, pods.node_uid, instance_type," +
		" req.name AS resource_name, req.quantity,\n    100 * req.quantity / sum(req.quantity) OVER w AS resource_percent\n" +
		"  FROM (\n    SELECT cluster_id, uid, node_uid\n    FROM pod_active\n    WHERE as_of <= '2019-12-02 10:00:00'\n" +
		"    UNION ALL -- We know these sets are disjoint.\n    SELECT cluster_id, uid, node_uid\n" +
		"    FROM pod_active_during\n    WHERE during @> '2019-12-02 10:00:00'::timestamp\n  ) pods\n" +
		"  NATURAL JOIN object\n  NATURAL JOIN namespaced_object\n  JOIN pod_resource_request req USING (cluster_id, uid)\n" +
		"  JOIN node_machine\n    ON node_machine.cluster_id = pods.cluster_id\n    AND node_machine.uid = pods.node_uid\n" +
		"  WHERE NOT EXISTS\n    (SELECT 1\n     FROM (\n       SELECT owner_uid AS uid\n       FROM owning_controller\n" +
		"       WHERE as_of <= '2019-12-02 10:00:00'\n         AND cluster_id = pods.cluster_id\n         AND uid = pods.uid\n" +
		"       UNION ALL -- We know these sets are disjoint.\n       SELECT owner_uid AS uid\n" +
		"       FROM owning_controller_during\n       WHERE during @> '2019-12-02 10:00:00'::timestamp\n" +
		"         AND cluster_id = pods.cluster_id\n         AND uid = pods.uid\n       LIMIT 1\n     ) owner\n" +
		"     JOIN object_type_meta\n       ON cluster_id = pods.cluster_id\n       AND object_type_meta.uid = owner.uid\n" +
		"       WHERE api_group = 'apps'\n         AND version = 'v1'\n         AND kind = 'DaemonSet'\n  )\n" +
		"  WINDOW w AS (PARTITION BY pods.cluster_id, node_uid, req.name)\n) by_resource\nJOIN object node_object\n" +
		"  ON node_object.cluster_id = by_resource.cluster_id\n  AND node_object.uid = by_resource.node_uid\n" +
		"LEFT OUTER JOIN node_pricing\n  ON node_pricing.cluster_id = by_resource.cluster_id\n" +
		"  AND node_pricing.uid = by_resource.node_uid\nLEFT OUTER JOIN node_topology\n" +
		"  ON node_topology.cluster_id = by_resource.cluster_id\n  AND node_topology.uid = by_resource.node_uid\n" +
		"LEFT OUTER JOIN (\n  SELECT cluster_id, uid, owner\n  FROM object_owner\n  WHERE as_of <= '2019-12-02 10:00:00'\n" +
		"  UNION ALL -- We know these sets are disjoint.\n  SELECT cluster_id, uid, owner\n  FROM object_owner_during\n" +
		"  WHERE during @> '2019-12-02 10:00:00'::timestamp\n) pod_owners\n  ON pod_owners.cluster_id = by_resource.cluster_id\n" +
		"  AND pod_owners.uid = by_resource.uid\nGROUP BY by_resource.cluster_id, node_name, instance_type, pricing_model, " +
		"region, zone, namespace, pod_name, owner;"
	tableNames = GetTableNamesInQuery(query)
	assert.Len(t, tableNames, 13)
}

func TestUilts_GetTidySQL(t *testing.T) {
	// TODO: syntax error for catalog.sampledb.abc from sqlparser
	// Fix https://github.com/uber/athenadriver/issues/5 when getting a chance
	assert.Equal(t, GetTidySQL(" SELECT * from catalog.sampledb.abc "), "SELECT * from catalog.sampledb.abc")
	assert.Equal(t, GetTidySQL(""), "")
	assert.Equal(t, GetTidySQL("DESC abc"), "DESC abc")
	assert.Equal(t, GetTidySQL("TRUNCATE abc"), "truncate table abc")
	assert.Equal(t, GetTidySQL("select"), "select")
	assert.Equal(t, GetTidySQL("drop table abc "), "drop table abc")
	assert.Equal(t, GetTidySQL("/**/ "), "")
	assert.Equal(t, GetTidySQL("/* select 1 */ select 1;"), "select 1")
	assert.Equal(t, GetTidySQL("SHOW FUNCTIONS;"), "show FUNCTIONS")
	assert.Equal(t, GetTidySQL("SELECT 1"), "select 1")
	assert.Equal(t, GetTidySQL("DROP TABLE ABC "), "drop table ABC")
	assert.Equal(t, GetTidySQL("/**/ "), "")
	assert.Equal(t, GetTidySQL("/* SELECT 1 */ SELECT 1;"), "select 1")
	assert.Equal(t, GetTidySQL("/* SELECT 1 */ SELECT 1 from;"), "SELECT 1 from;")
	assert.Equal(t, GetTidySQL(" select  *  from catalog.sampledb.abc "), "select  *  from catalog.sampledb.abc")
	assert.Equal(t, GetTidySQL(" select \"$path\" from sampledb.abc "), "select \"$path\" from sampledb.abc")
}

func TestUilts_GetCost(t *testing.T) {
	assert.Equal(t, getCost(0), 0.0)
	assert.Equal(t, getCost(1), getPrice10MB())
	assert.Equal(t, getCost(10*1024*1024*13), getPriceOneByte()*10*1024*1024*13)
}

func TestUtils_IsQID(t *testing.T) {
	assert.False(t, IsQID(`select "a44f8e61-4cbb-429a-b7ab-bea2c4a5caed"`))
	assert.True(t, IsQID("a44f8e61-4cbb-429a-b7ab-bea2c4a5caed"))
	assert.False(t, IsQID("a44f8e61-4cbb-429a-b7ab-bea2c4a5caeD"))
	assert.False(t, IsQID("a44f8e61"))
}

func Test_newHeaderResultPage(t *testing.T) {
	colName := "_col0"
	qid := "123"
	columnNames := []*string{&colName}
	columnTypes := []string{"string"}
	data := make([][]*string, 1)
	data[0] = []*string{&qid}
	page := newHeaderResultPage(columnNames, columnTypes, data)
	assert.NotNil(t, page)
}
