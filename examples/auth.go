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
	"os"

	secret "github.com/uber/athenadriver/examples/constants"
	drv "github.com/uber/athenadriver/go"
)

// To use athenadriver's Config for authentication
func useAthenaDriverConfigForAuth() {
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
	println("with AthenaDriver Config:", i)
}

// AWS_SDK_LOAD_CONFIG is used here for multiple use cases
// - use AWS CLI's Config for authentication
// - use in AWS Lambda where access ID and key are not required
// - assume role where access ID and key are not required
// Ref: https://github.com/uber/athenadriver/pull/10
func useAWSCLIConfigForAuth() {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	// 1. Set AWS Credential in Driver Config.
	conf, err := drv.NewDefaultConfig(secret.OutputBucketProd, drv.DefaultRegion,
		drv.DummyAccessID, drv.DummySecretAccessKey)
	if err != nil {
		println(err.Error())
		return
	}
	// 2. Open Connection.
	db, err := sql.Open(drv.DriverName, conf.Stringify())
	if err != nil {
		println(err.Error())
		return
	}
	// 3. Query and print results
	var i int
	err = db.QueryRow("SELECT 456").Scan(&i)
	if err != nil {
		println(err.Error())
	}
	println("with AWS CLI Config:", i)
	os.Unsetenv("AWS_SDK_LOAD_CONFIG")
}

// To use AWS CLI's Config for authentication with non-default profile set up by env variable AWS_PROFILE
// Refer: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
func useAWSCLIConfigForAuthProfileByEnv(profile string) {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	os.Setenv("AWS_PROFILE", profile)
	// 1. Set AWS Credential in Driver Config.
	conf, err := drv.NewDefaultConfig(secret.OutputBucketDev, drv.DummyRegion,
		drv.DummyAccessID, drv.DummySecretAccessKey)
	if err != nil {
		return
	}
	// 2. Open Connection.
	db, _ := sql.Open(drv.DriverName, conf.Stringify())
	// 3. Query and print results
	var i int
	_ = db.QueryRow("SELECT 789").Scan(&i)
	println("with AWS CLI Config With Profile:", i)
	os.Unsetenv("AWS_SDK_LOAD_CONFIG")
}

// To use AWS CLI's Config for authentication with a manually set up non-default profile
// Refer: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
func useAWSCLIConfigForAuthProfileByManualSetup(profile string) {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	// 1. Set AWS Credential in Driver Config.
	conf, err := drv.NewDefaultConfig(secret.OutputBucketDev, drv.DummyRegion,
		drv.DummyAccessID, drv.DummySecretAccessKey)
	if err != nil {
		return
	}
	conf.SetAWSProfile(profile)
	// 2. Open Connection.
	db, _ := sql.Open(drv.DriverName, conf.Stringify())
	// 3. Query and print results
	var i int
	_ = db.QueryRow("SELECT 789").Scan(&i)
	println("with AWS CLI Config With Profile:", i)
	os.Unsetenv("AWS_SDK_LOAD_CONFIG")
}

func main() {
	useAthenaDriverConfigForAuth()
	useAWSCLIConfigForAuth()
	useAWSCLIConfigForAuthProfileByEnv("henry")
	useAWSCLIConfigForAuthProfileByManualSetup("profile-development")
}

/*
Sample Output:
with AthenaDriver Config: 123
with AWS CLI Config: 456
with AWS CLI Config With Profile: 789
*/
