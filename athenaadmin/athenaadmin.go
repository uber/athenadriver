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
	"flag"
	"fmt"
	drv "github.com/uber/athenadriver/go"
)

var QID = flag.String("q", "1234", "the query id to be removed")
var table = flag.String("table", "henry.default.dummy", "the table name with format account.db_name.table_name")
var lastmodified = flag.Int64("lastmodified", 0, "the table's last modified time(int64)")
var s3location = flag.String("location", "", "the table's S3 location")
var report = flag.Bool("report", false, "to generate user report")

func removeQID(){
	cacheClient := drv.GetCacheInRedis()
	err := cacheClient.RemoveQID(*QID)
	if err != nil {
		println(err.Error())
	}
}

func updateTableLastModified(){
	cacheClient := drv.GetCacheInRedis()
	cacheClient.SetTableLastModified(*table, *lastmodified)
}

// https://docs.google.com/document/d/1hf6IzerIIEY0Xd9e7tUuT7Sa0ROq0L4YODlbGBh7S9A/edit#bookmark=id.6oxz4x3t955h
func updateTableS3Location(){
	if *s3location == "" {
		println("please provide a valid S3 location")
		return
	}
	cacheClient := drv.GetCacheInRedis()
	cacheClient.SetTableS3Location(*table, *s3location)
}

func generateUserReport(){
	if *report {
		cacheClient := drv.GetCacheInRedis()
		users := cacheClient.GetAllUsers()
		for _, u := range users {
			cacheClient.PrintStatsForUser(u)
			println("========================= Last 10 Queries =============================")
			QIDDateTime := cacheClient.GetQueryLogOfUser(u)
			for i, qt := range QIDDateTime  {
				q := cacheClient.GetQuery(qt[0])
				fmt.Printf("%d,%s,%s,%s\n", i+1, qt[0], q, qt[1])
			}
		}
	}
}

// To manage cache server, to generate report
func main() {
	flag.Parse()
	if *table == "henry.default.dummy" {
		println("please provide table name in format aws_account.db_name.table_name")
		return
	}
	if *report {
		generateUserReport()
		return
	}
	if *lastmodified !=0 {
		updateTableLastModified()
	}
	if *s3location != "" {
		updateTableS3Location()
	}

}
