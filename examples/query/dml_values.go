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

	secret "github.com/uber/athenadriver/examples/constants"
	drv "github.com/uber/athenadriver/go"
	"go.uber.org/zap"
)

// main will query Athena and print all columns and rows information in csv format
func main() {
	// 1. Set AWS Credential in Driver Config.
	conf, err := drv.NewDefaultConfig(secret.OutputBucket, secret.Region,
		secret.AccessID, secret.SecretAccessKey)
	if err != nil {
		panic(err)
		return
	}
	// 2. Open Connection.
	dsn := conf.Stringify()
	db, _ := sql.Open(drv.DriverName, dsn)

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	// 3. Query cancellation after 2 seconds
	ctx := context.WithValue(context.Background(), drv.LoggerKey, logger)
	rows, err := db.QueryContext(ctx, "values 1,2,3")
	if err != nil {
		println(err.Error())
		return
	}
	defer rows.Close()
	println(drv.ColsRowsToCSV(rows))
	rows, err = db.QueryContext(ctx, "VALUES\n    (1, 'a'),\n    (2, 'b'),\n    (3, 'c')")
	if err != nil {
		println(err.Error())
		return
	}
	defer rows.Close()
	println(drv.ColsRowsToCSV(rows))

}

/*
Sample output:
_col0
1
2
3

_col0,_col1
1,a
2,b
3,c
*/

/*
line 1:1: mismatched input 'call' expecting
{'(', 'select', 'from', 'add', 'desc', 'with', 'values', 'create', 'table',
	'insert', 'delete', 'describe', 'explain', 'show', 'use', 'drop', 'alter',
	'map', 'set', 'reset', 'start', 'commit', 'rollback', 'reduce', 'refresh',
	'clear', 'cache', 'uncache', 'dfs', 'truncate', 'analyze', 'list', 'revoke',
	'grant', 'lock', 'unlock', 'msck', 'export', 'import', 'load'}
(service: amazonathena; status code: 400; error code: invalidrequestexception; request id: 2f85b55c-4117-4ad2-9cd5-32d9777574b7)


line 1:1: extraneous input 'map' expecting
{'(', 'select', 'desc', 'using', 'with', 'values', 'create', 'table', 'insert',
'delete', 'describe', 'grant', 'revoke', 'explain', 'show', 'use', 'drop', 'alter',
'set', 'reset', 'start', 'commit', 'rollback', 'call', 'prepare', 'deallocate', 'execute'}
(service: amazonathena; status code: 400; error code: invalidrequestexception; request id: 5a343890-e5c8-4a4c-835a-0d6c0cd8d650)
*/
