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
	"database/sql"
	"log"

	secret "github.com/uber/athenadriver/examples/constants"

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

	wgTags := drv.NewWGTags()
	wgTags.AddTag("Uber User", "henry.wu")
	wgTags.AddTag("Uber ID", "123456")
	wgTags.AddTag("Uber Role", "SDE")
	// Specify workgroup name henry_wu should be used for the following query
	wg := drv.NewWG("henry_wu", nil, wgTags)
	_ = conf.SetWorkGroup(wg)
	// comment out the line below to allow remote workgroup creation and the query will be successful!!!
	//conf.SetWGRemoteCreationAllowed(false)

	// 2. Open Connection.
	dsn := conf.Stringify()
	db, _ := sql.Open(drv.DriverName, dsn)
	// 3. Query and print results
	rows, err := db.Query("select request_timestamp, url from sampledb.elb_logs limit 3")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var requestTimestamp string
	var url string
	for rows.Next() {
		if err := rows.Scan(&requestTimestamp, &url); err != nil {
			log.Fatal(err)
		}
		println(requestTimestamp + "," + url)
	}
}

/*
Sample Output:
2020/01/20 15:29:52 Workgroup henry_wu doesn't exist and workgroup remote creation is disabled.

After commenting out `conf.SetWGRemoteCreationAllowed(false)` at line 27:
2015-01-07T16:00:00.516940Z,https://www.example.com/articles/553
2015-01-07T16:00:00.902953Z,http://www.example.com/images/501
2015-01-07T16:00:01.206255Z,https://www.example.com/images/183

and you will see a new workgroup named `henry_wu` is created in AWS Athena console: https://us-east-2.console.aws.amazon.com/athena/workgroups/home?region=us-east-2
*/
