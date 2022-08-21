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

package main

import (
	"context"
	"database/sql"
	"log"

	secret "github.com/uber/athenadriver/examples/constants"
	"go.uber.org/zap"

	drv "github.com/uber/athenadriver/go"
)

func main() {
	// 1. Set AWS Credential in Driver Config.
	conf, err := drv.NewDefaultConfig(secret.OutputBucket, secret.Region,
		secret.AccessID, secret.SecretAccessKey)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 2. Open Connection.
	dsn := conf.Stringify()
	db, _ := sql.Open(drv.DriverName, dsn)
	// 3. Query and print results
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	var rows *sql.Rows
	rows, err = db.Query("DROP VIEW IF EXISTS sampledb.elb_logs_view;")
	if err != nil {
		log.Fatal(err)
		return
	}
	rows, err = db.Query("CREATE VIEW sampledb.elb_logs_view AS SELECT * FROM sampledb.elb_logs limit 1;")
	if err != nil {
		log.Println(err)
	}

	ctx := context.WithValue(context.Background(), drv.LoggerKey, logger)
	rows, err = db.QueryContext(ctx, "describe sampledb.elb_logs_view")
	if err != nil {
		log.Fatal(err)
		return
	}
	println(drv.ColsRowsToCSV(rows))
}

/*
Sample Output:
column,type
request_timestamp,varchar
elb_name,varchar
request_ip,varchar
request_port,integer
backend_ip,varchar
backend_port,integer
request_processing_time,double
backend_processing_time,double
client_response_time,double
elb_response_code,varchar
backend_response_code,varchar
received_bytes,bigint
sent_bytes,bigint
request_verb,varchar
url,varchar
protocol,varchar
user_agent,varchar
ssl_cipher,varchar
ssl_protocol,varchar

*/
