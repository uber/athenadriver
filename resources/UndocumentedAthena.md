# Athena Undocumented

Athena is a powerful query service. After digging into its official document and source code,
We found there are a lot of features and quirks undocumented or scattered in the internet here and there.
We did a lot of tests and debugging and wrote this article to share with Athena developers and 
Athena database users. We think this would help developers to understand Athena better to improve development productivity,
and help database users release the full power of AWS Athena.

## Data Type Supported

In Athena document website, there is a web page [Data Types Supported by Amazon Athena](https://docs.aws.amazon.com/athena/latest/ug/data-types.html). 
But according to our testing, what Athena supports are far more than that.

What Amazon did tell you is the following values are also valid and fully supported:

`json`, `varbinary`, `row`, `interval year to month`,  `interval day to second`,  `time`,  `time with time zone`,  `timestamp with time zone`

When querying against the above types, for the first three [`json`](https://github.com/uber/athenadriver/blob/master/examples/query/dml_select_json.go), [`varbinary`](https://github.com/uber/athenadriver/blob/master/examples/query/dml_select_geo.go), [`row`](https://github.com/uber/athenadriver/blob/master/examples/query/dml_select_row.go), **athenadriver** will return its string representation.
For the rest, a Go `time.Time` object will be returned.
 
In the following sample code, we use an SQL statement to `SELECT` som simple data of all the above types and then print them out.
At the same time, we enable the driver trace at debug level for this demonstration purpose.

```scala
package main

import (
	"context"
	"database/sql"
	secret "github.com/uber/athenadriver/examples/constants"
	drv "github.com/uber/athenadriver/go"
	"go.uber.org/zap"
)

func main() {
	// 1. Set AWS Credential in Driver Config.
	conf, err := drv.NewDefaultConfig("s3://henrywuqueryresults/",
		"us-east-2", "DummyAccessID", "DummySecretAccessKey")
	if err != nil {
		panic(err)
	}
	// 2. Open Connection.
	dsn := conf.Stringify()
	db, _ := sql.Open(drv.DBDriverName, dsn)
	// 3. Query and print results
	query := "SELECT JSON '\"Hello Athena\"', " +
		"ST_POINT(-74.006801, 40.70522), " +
		"ROW(1, 2.0),  INTERVAL '2' DAY, " +
		"INTERVAL '3' MONTH, " +
		"TIME '01:02:03.456', " +
		"TIME '01:02:03.456 America/Los_Angeles', " +
		"TIMESTAMP '2001-08-22 03:04:05.321 America/Los_Angeles';"
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	ctx := context.WithValue(context.Background(), drv.LoggerKey, logger)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	println(drv.ColsRowsToCSV(rows))
}
```

Sample output:

```bash
{"level":"debug","msg":"TM","athenaType":"json","goType":"string",
"str":"\"Hello Athena\""}
{"level":"debug","msg":"TM","athenaType":"varbinary","goType":"string",
"str":"00 00 00 00 01 01 00 00 00 20 25 76 6d 6f 80 52 c0 18 3e 22 a6 44 5a 44 40"}
{"level":"debug","msg":"TM","athenaType":"row","goType":"string",
"str":"{field0=1, field1=2.0}"}
{"level":"debug","msg":"TM","athenaType":"interval day to second","goType":"string",
"str":"2 00:00:00.000"}
{"level":"debug","msg":"TM","athenaType":"interval year to month","goType":"string",
"str":"0-3"}
{"level":"debug","msg":"TM","athenaType":"time","goType":"time.Time",
"str":"01:02:03.456"}
{"level":"debug","msg":"TM","athenaType":"time with time zone","goType":"time.Time",
"str":"01:02:03.456 America/Los_Angeles"}
{"level":"debug","msg":"TM","athenaType":"timestamp with time zone","goType":"time.Time",
"str":"2001-08-22 03:04:05.321 America/Los_Angeles"}
"Hello Athena",00 00 00 00 01 01 00 00 00 20 25 76 6d 6f 80 52 c0 18 3e 22 a6 44 5a 44 40,
{field0=1, field1=2.0},2 00:00:00.000,0-3,0000-01-01T01:02:03.456-07:52,
0000-01-01T01:02:03.456-07:52,2001-08-22T03:04:05.321-07:00
```

we can see `athenadriver` can handle all these undocumented types correctly.


## `ColumnInfo` has more number of cloumns than `Rows[0].Data`

> ![](pin.png)**Affected Statements: DESCRIBE TABLE/VIEW, SHOW SCHEMA/TABLE/...**

- Sample Query:

```sql
DESC sampledb.elb_logs
```

- Analysis:

![Column number mismatch issue example 1](issue_1.png)

We can see there are 3 columns according to `ColumnInfo` under `ResultSetMetadata`. But in the first row `Rows[0]`, we see there is only 1 field: `"elb_name \tstring    \t    "`. I would imagine there could have been 3 items in the `Data[0]`, but somehow the code author doesn't split it with tab(`\t`), so it ends up with only 1 item. The same issue happens for `SHOW` statement.

For more sample code, please check [util_desc_table.go](https://github.com/uber/athenadriver/blob/master/examples/query/util_desc_table.go), [util_desc_view.go](https://github.com/uber/athenadriver/blob/master/examples/query/util_desc_view.go), and [util_show.go](https://github.com/uber/athenadriver/blob/master/examples/query/util_show.go).

- `awsathendriver`'s Solution:

`athenadriver` fixes this issue by splitting `Rows[0].Data[0]` string with tab, and replace the original row with a new row which has the same number of data with columns.

## `ColumnInfo` has cloumns but `Rows` are empty

> ![](pin.png)**Affected Statements: [`CTAS`](https://docs.aws.amazon.com/athena/latest/ug/ctas.html), [CVAS](https://docs.aws.amazon.com/athena/latest/ug/views.html)\footnote{Create View as Select}, INSERT INTO**

Sample Query:

```sql
CREATE TABLE sampledb.elb_logs_copy WITH (
    format = 'TEXTFILE',
    external_location = 's3://external-location-henrywu/elb_logs_copy', 
    partitioned_by = ARRAY['ssl_protocol'])
AS SELECT * FROM sampledb.elb_logs
```

Analysis:

![Column number mismatch issue example 2](issue_3.png)

In the above [`CTAS`](https://docs.aws.amazon.com/athena/latest/ug/ctas.html) statement, we see there is one column of type `bigint` named
 `"rows"` in the resultset, but `ResultSet.Rows` is empty. Since there is no
  row, that one column doesn't make sense, or at least is confusing. The same
  issue happens for `INSERT INTO` statement.

- `awsathendriver`'s Solution:

Because this issue happens only in statements [`CTAS`](https://docs.aws.amazon.com/athena/latest/ug/ctas.html), [`CVAS`](https://docs.aws.amazon.com/athena/latest/ug/views.html), and `INSERT INTO
`, where `UpdateCount` is always valid and is the only meaningful information
 returned from Athena, `athenadriver` sets `UpdateCount` as the value of
  the returned row.

For more sample code, please check [ddl_ctas.go](https://github.com/uber/athenadriver/blob/master/examples/query/ddl_ctas.go), [ddl_cvas.go](https://github.com/uber/athenadriver/blob/master/examples/query/ddl_cvas.go), [dml_insert_into_select.go](https://github.com/uber/athenadriver/blob/master/examples/query/dml_insert_into_select.go) and [dml_insert_into_values.go](https://github.com/uber/athenadriver/blob/master/examples/query/dml_insert_into_values.go).


## When the Row resultset contains Header

[`GetQueryResults`](https://godoc.org/github.com/aws/aws-sdk-go/service/athena#Athena.GetQueryResults) and [`GetQueryResultsWithContext`](https://godoc.org/github.com/aws/aws-sdk-go/service/athena#Athena.GetQueryResultsWithContext) are the two functions to retrieve all Athena query results.

The resultset metadata includes both column and row details. If and Only if the statement is a select statement and the result set page is the first one, the first row is actually the column name strings, aka row header. Understanding this is very important to support all query statements, because when it is a SELECE statement, in order to make database driver behave consistently, we should skip the first row for the first page.


## What Query Types does Athena Support? 

According to [SQL Reference for Amazon Athena](https://docs.aws.amazon.com/athena/latest/ug/ddl-sql-reference.html), Amazon Athena supports a subset of Data Definition Language (DDL) and Data Manipulation Language (DML) statements, functions, operators, and data types. With some exceptions, Athena DDL is based on HiveQL DDL and Athena DML is based on Presto 0.172. But except the common DDL and DML statements, Athena also supports two UTILITY statements: `DESC` and `SHOW`.

In Athena source code, it is defined like:

```scala
	// The type of query statement that was run. DDL indicates DDL query statements.
	// DML indicates DML (Data Manipulation Language) query statements, such as
	// CREATE TABLE AS SELECT. UTILITY indicates query statements other than DDL
	// and DML, such as SHOW CREATE TABLE, or DESCRIBE <table>.
	StatementType *string `type:"string" enum:"StatementType"`
```

To ge the statement type, you can check [GetQueryExecutionOutput.QueryExecution.StatementType](https://docs.aws.amazon.com/athena/latest/APIReference/API_QueryExecution.html).

You can find all the statements' examples from [github.com/uber/athenadriver/examples/query](https://github.com/uber/athenadriver/tree/master/examples/query).

## How should we set `Database` in `athena.QueryExecutionContext{}`?

When queryng Athena, you can embed batabase string in the query when referring to the tables.
You can also provide the database name in addition to the query string as part of the function [StartQueryExecution](https://godoc.org/github.com/aws/aws-sdk-go/service/athena#Athena.StartQueryExecution)'s parameter like below:

```scala
resp, err := c.athenaAPI.StartQueryExecution(&athena.StartQueryExecutionInput{
	QueryString: aws.String(query),
	QueryExecutionContext: &athena.QueryExecutionContext{
		Database: aws.String(c.connector.AthenaConfig.GetDB()),
	},
	ResultConfiguration: &athena.ResultConfiguration{
		OutputLocation: aws.String(c.connector.AthenaConfig.GetOutputBucket()),
	},
	WorkGroup: aws.String(wg.Name),
})
```

However, the `Database` in line #4 bahaves inconsistently for different statements. It is very strange that for `SELECT`, `SHOW` statements, you can put any random string as Database, as long as there is database in query string, everything will work well. However for other statements like `DESC` , it will fail.

As this Athena use `default` as its default database name, but I figured out the best practice is actually setting it as `default`.


# Effective AWS Athena with `athenadriver` at Uber Technologies Inc

The following is an AWS Athena version of Scott Mayer's Effecitve C++.

## Does `athenadriver` support database reconnection?

Yes. `database/sql` maintains a connection pool internally and handles connection pooling, reconnecting, and retry logic for you.
One pitfall of writing Go sql application is cluttering the code with error-handling and retry.
I tested in my application with `athenadriver` by turning off and on Wifi and VPN, it works very well with database reconnection.

## Does `athenadriver` support batched query?
  
No. `athenadriver` is an implementation of `sql.driver` in Go `database/sql`, where there is no batch query support.
There might be some workaround for some specific case though. For instance, 
if you want to insert many rows, you can use [db.Exec](https://golang.org/pkg/database/sql/#DB.Exec) 
by replacing multiple inserts with one insert and multiple VALUES.
 
## How to get total row number of result set? Is there any way to randomly access row ?

You have to use `rows.Next()` to iterate all rows and use a counter to get row number. It is because Go `database/sql` was designed in a streaming query way with big data considered. That is why it only supports using `Next()` to iterate. So there is no way for random access of row. In Athena case, we only have random access of all the rows within one result page as the picture shown below:

![Encapsulation of driver.Rows in sql.Rows](sql_Rows.png) 

But due to encapsulation, more sepcifically the `rowsi` is _private_, we cannot access it directly like when we using Athena Go SDK. We have to use `Next()` to access it one by one.

## How to get the rows affected by my query?
  
The recommended way is to use `DB.Exec()` to get it. Please refer to \ref{db-exec}.

You can get it with `DB.Query()` too. In the returned `ResultSet`, there is an `UpdateCount` member variable. If the query is one of [`CTAS`](https://docs.aws.amazon.com/athena/latest/ug/ctas.html), [`CVAS`](https://docs.aws.amazon.com/athena/latest/ug/views.html) and `INSERT INTO`, `UpdateCount` will contain meaningful value. The result will be of a one row and one column. The column name is `rows`, and the row is an `int`, which is exactly `UpdateCount`. I would suggest to use `QueryRow` or `QueryRowContext` since it is a one-row result. By the way, the document for [`GetQueryResults`](https://docs.aws.amazon.com/athena/latest/APIReference/API_GetQueryResults.html) seems not very accurate.

![UpdateCount for CTAS, VTAS, and INSERT INTO](issue_2.png)

In practice, not only [`CTAS`](https://docs.aws.amazon.com/athena/latest/ug/ctas.html) statement but also `CVAS` and `INSERT INTO` will make a meaningful `UpdateCount`.

## How to makes database/table/view more discoverable? (from @allenw)

There are three ways.

###  `DESCRIBE` or `DESC` statment

The first is to use `DESC` statement. If you want to get details of column definition, you can just run something like:

```sql
DESC sampledb.elb_logs
```
Sample output(truncated for demo purpose):

| col_name                | data_type | comment | 
|-------------------------|-----------|---------| 
| request_timestamp       | string    |         | 
| elb_name                | string    |         | 
| request_ip              | string    |         | 


### `SELECT * FROM information_schema`

Athena provides `information_schema` schema so the second way is to query it in Athena.

Sample Query:

```sql
SELECT * FROM information_schema.columns
	where table_schema='sampledb' and table_name='elb_logs';
```
Sample Output(truncated for demo purpose):

| table_catalog  | table_schema | table_name | column_name             | ordinal_position | data_type |
|----------------|--------------|------------|-------------------------|------------------|-----------|
| awsdatacatalog | sampledb     | elb_logs   | request_timestamp       | 1                | varchar   |
| awsdatacatalog | sampledb     | elb_logs   | elb_name                | 2                | varchar   |
| awsdatacatalog | sampledb     | elb_logs   | request_ip              | 3                | varchar   |


### `SHOW`

The third way is to use `SHOW` statement.

Sample Query:

```sql
SHOW TBLPROPERTIES sampledb.elb_logs
```

Sample Output:

| prpt_name             | prpt_value | 
|-----------------------|------------| 
| EXTERNAL              | TRUE       | 
| transient_lastDdlTime | 1480278335 | 

Sample Query:
  
```sql
SHOW CREATE TABLE `sampledb.elb_logs`
```
Sample Output(reformatted for demo purpose):

```sql
createtab_stmt
CREATE EXTERNAL TABLE `sampledb.elb_logs`(
`request_timestamp` string COMMENT '', 
`elb_name` string COMMENT '', 
...
`ssl_protocol` string COMMENT '')
ROW FORMAT SERDE 
'org.apache.hadoop.hive.serde2.RegexSerDe' 
WITH SERDEPROPERTIES ( 
'input.regex'='([^ ]*) ([^ ]*) ([^ ]*):([0-9]*) ([^ ]*):([0-9]*) ([.0-9]*)
 ([.0-9]*) ([.0-9]*) (-|[0-9]*) (-|[0-9]*) ([-0-9]*) ([-0-9]*) \\\"([^ ]*)
 ([^ ]*) (- |[^ ]*)\\\" (\"[^\"]*\") ([A-Z0-9-]+) ([A-Za-z0-9.-]*)$') 
STORED AS INPUTFORMAT 
'org.apache.hadoop.mapred.TextInputFormat' 
OUTPUTFORMAT 
'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION
's3://athena-examples-us-east-2/elb/plaintext'
TBLPROPERTIES (
'transient_lastDdlTime'='1480278335')
```
More details are available at [DDL Statements](https://docs.aws.amazon.com/athena/latest/ug/language-reference.html).

## Does Athena allow to `DELETE` rows?

No. `DELETE` is not supported even in Athena database engine level, but `DROP` table and view are supported.

## Does Athena allow to `UPDATE` rows?

No. `UPDATE` is not supported in Athena database engine level, but `ALTER` table and database are allowed.

## Does Athena support cross account join or cross schema/database join? (from @allenw)

For cross schema/database join, it is as easy as in the normal SQL database.

Sample Queries:

```sql 
-- prepare two demo tables in different databases
CREATE TABLE default.elb_logs_new2 AS
 SELECT * FROM sampledb.elb_logs limit 1;
CREATE TABLE sampledb.elb_logs_new1 AS
 SELECT * FROM sampledb.elb_logs limit 10;
-- select-join query them
SELECT a.elb_name, a.url FROM elb_logs_new2 a 
 LEFT JOIN sampledb.elb_logs_new1 b ON a.request_timestamp=b.request_timestamp;
```

Sample Output:

| elb_name     | url                                  | 
|--------------|--------------------------------------| 
| elb_demo_009 | https://www.example.com/articles/746 |

Cross account join seems a malformed question on its own. It is impossible to join table in two different accounts because Athena and SQL don't provide a way to refer to the table. It won't work to refer to a able like  `account_1.database1.table1`. It is completely possible that two databases or tables in different accounts have the same name. There is no way to differentiate them in SQL statement.
An alternative is to create a new database or table under the same account and query them.
Because creating database in Athean is just create the metadata of some S3 data, it is a very lightweight operation.
So this is a feasible and recommended solution. 
If the s3 bucket owner and the account owner are different, you need to do [Cross-account Access](https://ocs.aws.amazon.com/athena/latest/ug/cross-account-permissions.html) setting.

## Contributing

As always, we welcome feedback and contributions. If you have any tips and findings about Athena, please feel free to contact [Henry Fuheng Wu](mailto:wufuheng@gmail.com).

