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

// The examples/perf folder contains stress, crash tests for performance and concurrency.
//
// How to prepare end-to-end crash test?
//
// 1. Prerequisites - AWS Credentials & S3 Query Result Bucket.
//
// To be able to query AWS Athena, you need to have an AWS account at Amazon AWS's website. To give it a shot,
// a free tier account is enough. You also need to have a pair of AWS access key ID and secret access key.
// You can get it from AWS Security Credentials section of Identity and Access Management (IAM).
// If you don't have one, please create it.
//
// In addition to AWS credentials, you also need an s3 bucket to store query result.
// Just go to AWS S3 web console page to create one. In the examples below,
// the s3 bucket I use is s3://henrywuqueryresults/.
//
// In most cases, you need the following 4 prerequisites :
//
//	S3 Output bucket
//	access key ID
//	secret access key
//	AWS region
//
// For more details on athenadriver's support on AWS credentials & S3 query result bucket,
// please refer to README section Support Multiple AWS Authorization Methods.
//
// 2. Installation athenadriver.
//
// Before Go 1.17, go get can be used to install athenadriver:
//
//   go get -u github.com/uber/athenadriver
// 
// Starting in Go 1.17, installing executables with go get is deprecated. go install may be used instead.
//
//   go install github.com/uber/athenadriver@latest
//
// 3. Integration Test.
//
// To Build it:
//
//	$cd $GOPATH/src/github.com/uber/athenadriver
//	$go build examples/perf/concurrency.go
//
// Run it and wait for some output and unplug your cable:
//
//	$ulimit -c unlimited
//	$sudo sysctl -w kernel.core_pattern=/tmp/core
//	$ GOTRACEBACK=crash ./concurrency > examples/perf/concurrency.output.`date +"%Y-%m-%d-%H-%M-%S"`.log
//	58,13,53,54,78,96,32,48,40,11,35,31,65,61,1,73,74,22,34,49,80,5,69,37,0,79,2020/02/09 13:49:29 error [38]RequestError: send request failed
//	caused by: Post https://athena.us-east-1.amazonaws.com/: dial tcp: lookup athena.us-east-1.amazonaws.com: no such host
//	2020/02/09 13:49:29 error [24]RequestError: send request failed
//	caused by: Post https://athena.us-east-1.amazonaws.com/: dial tcp: lookup athena.us-east-1.amazonaws.com: no such host
//	2020/02/09 13:49:29 error [64]RequestError: send request failed
//	caused by: Post https://athena.us-east-1.amazonaws.com/: dial tcp: lookup athena.us-east-1.amazonaws.com: no such host
//	2020/02/09 13:49:29 error [55]RequestError: send request failed
//	caused by: Post https://athena.us-east-1.amazonaws.com/: dial tcp: lookup athena.us-east-1.amazonaws.com: no such host
//	...
//	2020/02/09 13:49:29 error [95]RequestError: send request failed
//	caused by: Post https://athena.us-east-1.amazonaws.com/: dial tcp: lookup athena.us-east-1.amazonaws.com: no such host
//	2020/02/09 13:49:29 error [9]RequestError: send request failed
//	caused by: Post https://athena.us-east-1.amazonaws.com/: dial tcp: lookup athena.us-east-1.amazonaws.com: no such host
//	2020/02/09 13:49:29 error [89]RequestError: send request failed
//	caused by: Post https://athena.us-east-1.amazonaws.com/: dial tcp: lookup athena.us-east-1.amazonaws.com: no such host
//
// And now re-plugin your cable and wait for network coming back, you can see the program automatically reconnect, and resume to output correctly:
//
//	72,25,92,98,15,93,41,7,8,90,81,56,66,2,18,84,87,63,44,45,82,99,86,3,52,76,71,16,39,67,23,12,42,17,4,
package main
