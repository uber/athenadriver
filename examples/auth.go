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
	secret "github.com/uber/athenadriver/examples/constants"
	drv "github.com/uber/athenadriver/go"
	"os"
)

// To use athenadriver's Config for authentication
func useAthenaSQLConfigForAuth() {
	// 1. Set AWS Credential in Driver Config.
	conf, err := drv.NewDefaultConfig(secret.OutputBucketDev, secret.Region,
		secret.AccessID, secret.SecretAccessKey)
	if err != nil {
		return
	}
	// 2. Open Connection.
	db, _ := sql.Open(drv.DriverName, conf.Stringify())
	// 3. Query and print results
	var i int
	_ = db.QueryRow("SELECT 123").Scan(&i)
	println("with AthenaSQL Config:", i)
}

// To use AWS CLI's Config for authentication
func useAWSCLIConfigForAuth() {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	// 1. Set AWS Credential in Driver Config.
	conf, err := drv.NewDefaultConfig(secret.OutputBucketProd, drv.DummyRegion,
		drv.DummyAccessID, drv.DummySecretAccessKey)
	if err != nil {
		return
	}
	// 2. Open Connection.
	db, _ := sql.Open(drv.DriverName, conf.Stringify())
	// 3. Query and print results
	var i int
	_ = db.QueryRow("SELECT 456").Scan(&i)
	println("with AWS CLI Config:", i)
	os.Unsetenv("AWS_SDK_LOAD_CONFIG")
}

func main() {
	useAthenaSQLConfigForAuth()
	useAWSCLIConfigForAuth()
}

/*
Sample Output:
with AthenaSQL Config: 123
with AWS CLI Config: 456
*/
