# :shell: athenareader

`athenareader` is a utility tool which query S3 data via Athena from command line. It output query result in CSV format.

## Authentication Method

To avoid exposing access keys(Access Key ID and Secret Access Key) in command line, `athenareader` use [AWS CLI Config For Authentication](https://github.com/uber/athenadriver#use-aws-cli-config-for-authentication) method. Please make sure your environment variable **`AWS_SDK_LOAD_CONFIG`** is set.

## How to get/build/install `athenareader`

```
go get -u github.com/uber/athenadriver/athenareader
```

## How to use `athenareader`

```
$ athenareader --help
NAME
	athenareader - read athena data from command line

SYNOPSIS
	athenareader [-v] [-b output_bucket] [-d database_name] [-q query_string_or_file] [-r] [-a] [-m]

DESCRIPTION
  -a	Enable admin mode, so database write(create/drop) is allowed at athenadriver level
  -b string
    	Athena resultset output bucket (default "s3://query-results-bucket-henrywu/")
  -d string
    	The database you want to query (default "default")
  -m	Enable moneywise mode to display the query cost as the first line of the output
  -q string
    	The SQL query string or a file containing SQL string (default "select 1")
  -r	Display rows only, don't show the first row as columninfo
  -v	Print the current version and exit

EXAMPLES

	$ athenareader -d sampledb -q "select request_timestamp,elb_name from elb_logs limit 2"
	request_timestamp,elb_name
	2015-01-03T00:00:00.516940Z,elb_demo_004
	2015-01-03T00:00:00.902953Z,elb_demo_004

	$ athenareader -d sampledb -q "select request_timestamp,elb_name from elb_logs limit 2" -r
	2015-01-05T20:00:01.206255Z,elb_demo_002
	2015-01-05T20:00:01.612598Z,elb_demo_008

	$ athenareader -d sampledb -b s3://my-athena-query-result -q tools/query.sql
	request_timestamp,elb_name
	2015-01-06T00:00:00.516940Z,elb_demo_009


	Add '-m' to enable moneywise mode. The first line will display query cost under moneywise mode.

	$ athenareader -b s3://athena-query-result -q 'select count(*) as cnt from sampledb.elb_logs' -m
	query cost: 0.00184898369752772851 USD
	cnt
	1356206


	Add '-a' to enable admin mode. Database write is enabled at driver level under admin mode.

	$ athenareader -b s3://athena-query-result -q 'DROP TABLE IF EXISTS depreacted_table' -a
	
AUTHOR
	Henry Fuheng Wu (henry.wu@uber.com)

REPORTING BUGS
	https://github.com/uber/athenadriver
```
