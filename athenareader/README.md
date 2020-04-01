# :shell: athenareader

`athenareader` is a utility tool which query S3 data via Athena from command line. It output query result in CSV format.

## Authentication Method

To avoid exposing access keys(Access Key ID and Secret Access Key) in command line, `athenareader` use [AWS CLI Config For Authentication](https://github.com/uber/athenadriver#use-aws-cli-config-for-authentication) method. Please make sure your environment variable **`AWS_SDK_LOAD_CONFIG`** is set.

## How to get/build/install `athenareader`

```
go get github.com/uber/athenadriver/athenareader
```

## How to use `athenareader`

```
athenareader --help
NAME
	athenareader - read athena data from command line

SYNOPSIS
	./athenareader [-b output_bucket] [-d database_name] [-q query_string_or_file] [-r]

DESCRIPTION
  -b string
    	Athena resultset output bucket (default "s3://query-results-bucket-henrywu/")
  -d string
    	The database you want to query (default "sampledb")
  -q string
    	The SQL query string or a file containing SQL string (default "select 1")
  -r	Display rows only, don't show the first row as columninfo

EXAMPLES

	$ athenareader -d sampledb -q "select request_timestamp,elb_name from elb_logs limit 2"
	request_timestamp,elb_name
	2015-01-03T00:00:00.516940Z,elb_demo_004
	2015-01-03T00:00:00.902953Z,elb_demo_004

	$ athenareader -d sampledb -q "select request_timestamp,elb_name from elb_logs limit 2" -r
	2015-01-05T20:00:01.206255Z,elb_demo_002
	2015-01-05T20:00:01.612598Z,elb_demo_008

	$ athenareader -d sampledb -q tools/query.sql
	request_timestamp,elb_name
	2015-01-06T00:00:00.516940Z,elb_demo_009

AUTHOR
	Henry Fuheng Wu(henry.wu@uber.com)

REPORTING BUGS
	https://github.com/uber/athenadriver
```
