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

package athenadriver

// TContextKey is a type for key in context.
type TContextKey string

const (
	// DriverName is the Name of this DB driver.
	DriverName = "awsathena"

	// DefaultBytesScannedCutoffPerQuery is 1G for every user.
	DefaultBytesScannedCutoffPerQuery = 1024 * 1024 * 1024

	// DefaultDBName is the default database name in Athena.
	DefaultDBName = "default"

	// DefaultWGName is the default workgroup name in Athena
	DefaultWGName = "primary"

	// DefaultRegion is the default region in Athena.
	DefaultRegion = "us-east-1"

	// TimestampUniXFormat is from https://docs.aws.amazon.com/athena/latest/ug/data-types.html.
	// https://stackoverflow.com/questions/20530327/origin-of-mon-jan-2-150405-mst-2006-in-golang
	// RFC3339 is not supported by AWS Athena. It uses session timezone!.
	TimestampUniXFormat = "2006-01-02 15:04:05.000"

	// ZeroDateTimeString is the invalid or zero result for a time.Time
	ZeroDateTimeString = "0001-01-01 00:00:00 +0000 UTC"

	// DateUniXFormat comes along the same way as TimestampUniXFormat.
	DateUniXFormat = "2006-01-02"

	// MetricsKey is the key for Metrics in context
	MetricsKey = TContextKey("MetricsKey")

	// LoggerKey is the key for Logger in context
	LoggerKey = TContextKey("LoggerKey")

	// DummyRegion is used when AWS CLI Config is used, ie AWS_SDK_LOAD_CONFIG is set
	DummyRegion = "dummy"

	// DummyAccessID is used when AWS CLI Config is used, ie AWS_SDK_LOAD_CONFIG is set
	DummyAccessID = "dummy"

	// DummySecretAccessKey is used when AWS CLI Config is used, ie AWS_SDK_LOAD_CONFIG is set
	DummySecretAccessKey = "dummy"
)

// https://docs.aws.amazon.com/athena/latest/ug/service-limits.html
const (
	// DDLQueryTimeout is DDL query timeout 600 minutes(unit second).
	DDLQueryTimeout = 600 * 60

	// DMLQueryTimeout is DML query timeout 30 minutes(unit second).
	DMLQueryTimeout = 30 * 60

	// PoolInterval is the interval between two status checks(unit second).
	PoolInterval = 3

	// The maximum allowed query string length is 262144 bytes,
	// where the strings are encoded in UTF-8.
	// This is not an adjustable quota. (unit bytes)
	MAXQueryStringLength = 262144
)

const digits01 = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
const digits10 = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999"

// AthenaColumnTypes is a fixed array of Athena Column Types. An array isn't immutable by nature; you can't make it constant.
var AthenaColumnTypes = [...]string{"tinyint", "smallint", "integer", "bigint", "float", "real", "double",
	"json", "char", "varchar", "varbinary", "row", "string", "binary",
	"struct", "interval year to month", "interval day to second", "decimal",
	"ipaddress", "array", "map", "unknown", "boolean", "date", "time", "time with time zone",
	"timestamp with time zone", "timestamp", "weird_type"}

// pseudo commands all start with `PC_`

// PCGetQID is the pseudo command of getting query execution id of an SQL
const PCGetQID = "get_query_id"

// PCGetQIDStatus is the pseudo command of getting status of a query execution id
const PCGetQIDStatus = "get_query_id_status"

// PCStopQID is the pseudo command to stop a query execution id
const PCStopQID = "stop_query_id"

// PCGetDriverVersion is the pseudo command to get the version of athenadriver
const PCGetDriverVersion = "get_driver_version"

// DriverVersion is athenadriver's version
const DriverVersion = "1.1.14"
