// Copyright (c) 2020 Uber Technologies, Inc.
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

package main

import (
	"database/sql"
	"flag"
	secret "github.com/uber/athenadriver/examples/constants"
	drv "github.com/uber/athenadriver/go"
	"io/ioutil"
	"log"
	"fmt"
	"os"
)

var CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

// main will query Athena and print all columns and rows information in csv format
func main() {
	var bucket = flag.String("b", secret.OutputBucket, "Athena resultset output bucket")
	var database = flag.String("d", "sampledb", "The database you want to query")
	var query = flag.String("q", "select 1", "The SQL query string or a file containing SQL string")
	var rowOnly = flag.Bool("r", false, "Display rows only, don't show the first row as columninfo")

	flag.Usage = func() {
		pre_body:="NAME\n\tathenareader - read athena data from command line\n\n"
		desc := "\nEXAMPLES\n\n" +
			"\t$ athenareader -d sampledb -q \"select request_timestamp,elb_name from elb_logs limit 2\"\n"+
			"\trequest_timestamp,elb_name\n" +
			"\t2015-01-03T00:00:00.516940Z,elb_demo_004\n" +
			"\t2015-01-03T00:00:00.902953Z,elb_demo_004\n\n" +
			"\t$ athenareader -d sampledb -q \"select request_timestamp,elb_name from elb_logs limit 2\" -r\n" +
			"\t2015-01-05T20:00:01.206255Z,elb_demo_002\n" +
			"\t2015-01-05T20:00:01.612598Z,elb_demo_008\n\n" +
			"\t$ athenareader -d sampledb -q tools/query.sql\n" +
			"\trequest_timestamp,elb_name\n" +
			"\t2015-01-06T00:00:00.516940Z,elb_demo_009\n\n" +
			"AUTHOR\n\tHenry Fuheng Wu(henry.wu@uber.com)\n\n" +
			"REPORTING BUGS\n\thttps://github.com/uber/athenadriver\n"
		fmt.Fprintf(CommandLine.Output(), pre_body)
		fmt.Fprintf(CommandLine.Output(),
			"SYNOPSIS\n\t%s [-b output_bucket] [-d database_name] [-q query_string_or_query_file] [-r]\n\nDESCRIPTION\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(CommandLine.Output(), desc)
	}

	flag.Parse()
	// 1. Set AWS Credential in Driver Config.
	conf, err := drv.NewDefaultConfig(*bucket, secret.Region, secret.AccessID, secret.SecretAccessKey)
	conf.SetDB(*database)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 2. Open Connection.
	dsn := conf.Stringify()
	db, _ := sql.Open(drv.DriverName, dsn)
	// 3. Query and print results
	sqlString := *query
	if _, err := os.Stat(*query); err == nil {
		b, err := ioutil.ReadFile(*query)
		if err != nil {
			fmt.Print(err)
		}
		sqlString = string(b) // convert content to a 'string'
	}
	rows, err := db.Query(sqlString)
	if err != nil {
		println(err.Error())
		return
	}
	defer rows.Close()
	if *rowOnly {
		println(drv.RowsToCSV(rows))
		return
	}
	println(drv.ColsRowsToCSV(rows))
}
