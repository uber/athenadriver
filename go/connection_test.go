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

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"math/rand"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/stretchr/testify/assert"
)

var regions = []string{"ap-east-1", "eu-central-1", "eu-north-1", "eu-west-1", "eu-west-2", "eu-west-3",
	"me-south-1", "us-east-1", "us-west-1", "ap-northeast-1", "ap-northeast-2", "ap-southeast-1",
	"ca-central-1", "us-east-2", "ap-south-1", "ap-southeast-2", "us-west-2",
}

func TestReadOnlyCTAS1(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetReadOnly(true)
	dsn := testConf.Stringify()
	db, _ := sql.Open(DriverName, dsn)
	_, err := db.QueryContext(context.Background(), "CREATE TABLE sampledb."+
		"elb_logs_new AS "+
		"SELECT * FROM sampledb.elb_logs limit 10;")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "writing to Athena database is disallowed in read-only mode")

	query := randString(MAXQueryStringLength*10) + "?"
	args := []driver.Value{query}
	r, err := db.ExecContext(context.Background(), query, args)
	assert.NotNil(t, err)
	assert.Equal(t, r, nil)

	r, err = db.ExecContext(context.Background(), query, "")
	assert.NotNil(t, err)
	assert.Equal(t, r, nil)

	r, err = db.ExecContext(context.Background(), query)
	assert.NotNil(t, err)
	assert.Equal(t, r, nil)

	query = "?"
	r, err = db.ExecContext(context.Background(), query, "")
	assert.NotNil(t, err)
	assert.Equal(t, r, nil)

	e := db.Ping()
	assert.NotNil(t, e)
}

func TestReadOnlyCTAS2(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetReadOnly(true)
	dsn := testConf.Stringify()
	db, _ := sql.Open(DriverName, dsn)
	_, err := db.QueryContext(context.Background(), " CREATE TABLE sampledb."+
		"elb_logs_new AS "+
		"SELECT * FROM sampledb.elb_logs limit 10;")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "writing to Athena database is disallowed in read-only mode")
}

func TestReadOnlyCTAS3(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetReadOnly(true)
	dsn := testConf.Stringify()
	db, _ := sql.Open(DriverName, dsn)
	_, err := db.QueryContext(context.Background(), " cReate TABLE sampledb."+
		"elb_logs_new AS "+
		"SELECT * FROM sampledb.elb_logs limit 10;")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "writing to Athena database is disallowed in read-only mode")
}

func TestReadOnlyDROP(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetReadOnly(true)
	dsn := testConf.Stringify()
	db, _ := sql.Open(DriverName, dsn)
	_, err := db.QueryContext(context.Background(), " drop table test")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "writing to Athena database is disallowed in read-only mode")
}

func TestConnection_Prepare(t *testing.T) {
	testConf := NewNoOpsConfig()
	connector := &SQLConnector{
		config: testConf,
		tracer: NewDefaultObservability(testConf),
	}

	conn, err := connector.Connect(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	prepStatement, err := conn.Prepare("select 123")
	assert.NotNil(t, prepStatement)
	assert.Nil(t, err)

	query := randString(MAXQueryStringLength * 10)
	prepStatement, err = conn.Prepare(query)
	assert.NotNil(t, err)
	assert.Nil(t, prepStatement)
}

func TestConnection_Begin(t *testing.T) {
	testConf := NewNoOpsConfig()
	connector := &SQLConnector{
		config: testConf,
		tracer: NewDefaultObservability(testConf),
	}

	conn, err := connector.Connect(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	transaction, err := conn.Begin()
	assert.Nil(t, transaction)
	assert.Equal(t, err.Error(), "Athena doesn't support transaction statements")
}

func TestConnection_Close(t *testing.T) {
	testConf := NewNoOpsConfig()
	connector := &SQLConnector{
		config: testConf,
		tracer: NewDefaultObservability(testConf),
	}

	conn, err := connector.Connect(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	assert.Nil(t, conn.Close())

}

func TestConnection_QueryContext(t *testing.T) {
	testConf := NewNoOpsConfig()
	connector := &SQLConnector{
		config: testConf,
		tracer: NewDefaultObservability(testConf),
	}

	conn, err := connector.Connect(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, conn)
}

func TestConnection_BeginTx(t *testing.T) {
	c := createTestConnection(t)
	value := driver.NamedValue{Value: uint64(0)}
	err := c.CheckNamedValue(&value)
	if err != nil {
		t.Fatal("uint64 not convertible", err)
	}
	tx, err := c.BeginTx(context.Background(),
		&sql.TxOptions{Isolation: sql.
			LevelSerializable})
	assert.Nil(t, tx)
	assert.Equal(t, err.Error(), "Athena doesn't support transaction statements")
}

func TestConnection_Transaction(t *testing.T) {
	db, _ := sql.Open(DriverName, NewNoOpsConfig().Stringify())
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	assert.Nil(t, tx)
	assert.Equal(t, err.Error(), "sql: driver does not support non-default isolation level")
}

func TestConnection_InterpolateParams(t *testing.T) {
	c := createTestConnection(t)
	q, err := c.interpolateParams("SELECT ?+?", []driver.Value{int64(42), "gopher"})
	if err != nil {
		t.Errorf("Expected err=nil, got %#v", err)
		return
	}
	expected := `SELECT 42+'gopher'`
	if q != expected {
		t.Errorf("Expected: %q\nGot: %q", expected, q)
	}
}

func TestInterpolateParamsTooManyPlaceholders(t *testing.T) {
	c := createTestConnection(t)
	q, err := c.interpolateParams("SELECT ?+?", []driver.Value{int64(42)})
	if err != ErrInvalidQuery {
		t.Errorf("Expected err=ErrInvalidQuery, got err=%#v, q=%#v", err, q)
	}
}

func TestConnection_InterpolateParams_Query(t *testing.T) {
	c := createTestConnection(t)
	query := randString(MAXQueryStringLength*10) + "?"
	q, err := c.interpolateParams(query, []driver.Value{query})
	assert.Equal(t, q, "")
	assert.NotNil(t, err)
}

func TestConnection_InterpolateParams_Query2(t *testing.T) {
	c := createTestConnection(t)
	q, err := c.interpolateParams("?", []driver.Value{aType{S: "abc"}})
	assert.Equal(t, q, "")
	assert.NotNil(t, err)

	q, err = c.interpolateParams("?", []driver.Value{1})
	assert.Equal(t, q, "")
	assert.NotNil(t, err)
}

func TestConnection_InterpolateParams_Bool(t *testing.T) {
	c := createTestConnection(t)
	q, err := c.interpolateParams("?", []driver.Value{true})
	assert.Equal(t, q, "1")
	assert.Nil(t, err)
	q, err = c.interpolateParams("?", []driver.Value{false})
	assert.Equal(t, q, "0")
	assert.Nil(t, err)
	q, err = c.interpolateParams("?", []driver.Value{int64(1)})
	assert.Equal(t, q, "1")
	assert.Nil(t, err)
	q, err = c.interpolateParams("?", []driver.Value{uint64(1)})
	assert.Equal(t, q, "1")
	assert.Nil(t, err)
	q, err = c.interpolateParams("?", []driver.Value{float64(1.1)})
	assert.Equal(t, q, "1.1")
	assert.Nil(t, err)
	q, err = c.interpolateParams("?", []driver.Value{time.Time{}})
	assert.Equal(t, q, "'0000-00-00'")
	assert.Nil(t, err)
	q, err = c.interpolateParams("?", []driver.Value{time.Now()})
	assert.NotEqual(t, q, "'0000-00-00'")
	assert.Nil(t, err)
	q, err = c.interpolateParams("?", []driver.Value{[]byte{'0'}})
	assert.Equal(t, q, "_binary'0'")
	assert.Nil(t, err)
	q, err = c.interpolateParams("?", []driver.Value{nil})
	assert.Equal(t, q, "NULL")
	assert.Nil(t, err)
	q, err = c.interpolateParams("123?4", []driver.Value{nil})
	assert.Equal(t, q, "123NULL4")
	assert.Nil(t, err)
	q, err = c.interpolateParams("?", []driver.Value{time.Time{}.Add(1 * time.Nanosecond)})
	assert.Equal(t, q, "'0001-01-01 00:00:00'")
	assert.Nil(t, err)
}

// We don't support placeholder in string literal for now.
// https://github.com/go-sql-driver/mysql/pull/490
func TestInterpolateParamsPlaceholderInString(t *testing.T) {
	c := createTestConnection(t)

	q, err := c.interpolateParams("SELECT 'abc?xyz',?", []driver.Value{int64(42)})
	// When InterpolateParams support string literal, this should return `"SELECT 'abc?xyz', 42`
	if err != ErrInvalidQuery {
		t.Errorf("Expected err=ErrInvalidQuery, got err=%#v, q=%#v", err, q)
	}
}

func TestInterpolateParamsUint64(t *testing.T) {
	c := createTestConnection(t)

	q, err := c.interpolateParams("SELECT ?", []driver.Value{uint64(42)})
	assert.Nil(t, err)
	assert.Equal(t, q, "SELECT 42")
}

func TestCheckNamedValue(t *testing.T) {
	c := createTestConnection(t)
	value := driver.NamedValue{Value: uint64(0)}
	err := c.CheckNamedValue(&value)
	assert.Nil(t, err)
	assert.Equal(t, value.Value, int64(0))
}

func createTestConnection(t *testing.T) *Connection {
	t.Parallel()
	testConf := NewNoOpsConfig()
	staticCredentials := credentials.NewStaticCredentials(testConf.GetAccessID(),
		testConf.GetSecretAccessKey(),
		testConf.GetSessionToken())
	awsConfig := &aws.Config{
		Region:      aws.String(testConf.GetRegion()),
		Credentials: staticCredentials,
	}
	awsAthenaSession, err := session.NewSession(awsConfig)
	assert.Nil(t, err)
	athenaAPI := athena.New(awsAthenaSession)
	c := &Connection{
		athenaAPI: athenaAPI,
		connector: NoopsSQLConnector(),
	}
	return c
}

func TestConnection_QueryContext2(t *testing.T) {
	t.Parallel()
	c := &Connection{
		athenaAPI: newMockAthenaClient(),
		connector: NoopsSQLConnector(),
	}
	driverRows, err := c.QueryContext(context.Background(), "StartQueryExecution_nil_error",
		[]driver.NamedValue{})
	assert.Nil(t, driverRows)
	assert.NotNil(t, err)

	driverRows, err = c.QueryContext(context.Background(), "StartQueryExecution_OK_GetQueryExecutionWithContext_QueryExecutionStateCancelled",
		[]driver.NamedValue{})
	assert.Nil(t, driverRows)
	assert.Equal(t, err, context.Canceled)

	driverRows, err = c.QueryContext(context.Background(), "StartQueryExecution_OK_GetQueryExecutionWithContext_QueryExecutionStateFailed",
		[]driver.NamedValue{})
	assert.Nil(t, driverRows)
	assert.Equal(t, err, ErrTestMockFailedByAthena)

}

func TestConnection_QueryContext3(t *testing.T) {
	t.Parallel()
	c := &Connection{
		athenaAPI: newMockAthenaClient(),
		connector: NoopsSQLConnector(),
	}
	var s3bucket string = "s3://query-results-henry-wu-us-east-2/"

	wgTags := NewWGTags()
	wgTags.AddTag("Uber User", "henry.wu")
	wgTags.AddTag("Uber Asset", "abc.efg")
	wg := NewDefaultWG("henry_wu", nil, wgTags)
	testConf := NewNoOpsConfig()
	err := testConf.SetOutputBucket(s3bucket)
	assert.Nil(t, err)
	err = testConf.SetRegion("us-east-1")
	assert.Nil(t, err)
	testConf.SetUser("henry.wu")
	testConf.SetDB("default") // default

	_ = testConf.SetWorkGroup(wg)
	c.connector.config = testConf
	driverRows, err := c.QueryContext(context.Background(), "StartQueryExecution_nil_error",
		[]driver.NamedValue{})
	assert.Nil(t, driverRows)
	assert.NotNil(t, err)
}

func TestConnection_QueryContext4(t *testing.T) {
	t.Parallel()
	nm := newMockAthenaClient()
	nm.GetWGStatus = true
	c := &Connection{
		athenaAPI: nm,
		connector: NoopsSQLConnector(),
	}
	var s3bucket string = "s3://query-results-henry-wu-us-east-2/"

	wgTags := NewWGTags()
	wgTags.AddTag("Uber User", "henry.wu")
	wgTags.AddTag("Uber Asset", "abc.efg")
	wg := NewDefaultWG("henry_wu", nil, wgTags)
	testConf := NewNoOpsConfig()
	err := testConf.SetOutputBucket(s3bucket)
	assert.Nil(t, err)
	err = testConf.SetRegion("us-east-1")
	assert.Nil(t, err)
	testConf.SetUser("henry.wu")
	testConf.SetDB("default") // default

	_ = testConf.SetWorkGroup(wg)
	c.connector.config = testConf

	driverRows, err := c.QueryContext(context.Background(), "StartQueryExecution_nil_error",
		[]driver.NamedValue{})
	assert.Nil(t, driverRows)
	assert.NotNil(t, err)
}

func TestConnection_QueryContext5(t *testing.T) {
	t.Parallel()
	nm := newMockAthenaClient()
	nm.GetWGStatus = true
	nm.WGDisabled = true

	c := &Connection{
		athenaAPI: nm,
		connector: NoopsSQLConnector(),
	}
	var s3bucket string = "s3://query-results-henry-wu-us-east-2/"

	wgTags := NewWGTags()
	wgTags.AddTag("Uber User", "henry.wu")
	wgTags.AddTag("Uber Asset", "abc.efg")
	wg := NewDefaultWG("henry_wu", nil, wgTags)
	testConf := NewNoOpsConfig()
	_ = testConf.SetOutputBucket(s3bucket)
	_ = testConf.SetRegion("us-east-1")
	testConf.SetUser("henry.wu")
	testConf.SetDB("default") // default
	_ = testConf.SetWorkGroup(wg)
	c.connector.config = testConf

	driverRows, err := c.QueryContext(context.Background(), "StartQueryExecution_nil_error",
		[]driver.NamedValue{})
	assert.Nil(t, driverRows)
	assert.NotNil(t, err)
}

func TestConnection_QueryContext6(t *testing.T) {
	t.Parallel()
	nm := newMockAthenaClient()
	c := &Connection{
		athenaAPI: nm,
		connector: NoopsSQLConnector(),
	}
	var s3bucket string = "s3://query-results-henry-wu-us-east-2/"

	wgTags := NewWGTags()
	wgTags.AddTag("Uber User", "henry.wu")
	wgTags.AddTag("Uber Asset", "abc.efg")
	wg := NewDefaultWG("henry_wu", nil, wgTags)
	testConf := NewNoOpsConfig()
	err := testConf.SetOutputBucket(s3bucket)
	assert.Nil(t, err)
	err = testConf.SetRegion("us-east-1")
	assert.Nil(t, err)
	testConf.SetUser("henry.wu")
	testConf.SetDB("default") // default
	testConf.SetWGRemoteCreationAllowed(false)

	_ = testConf.SetWorkGroup(wg)
	c.connector.config = testConf

	e := c.Ping(context.Background())
	assert.Equal(t, e, driver.ErrBadConn)
	driverRows, err := c.QueryContext(context.Background(), "StartQueryExecution_nil_error",
		[]driver.NamedValue{})
	assert.Nil(t, driverRows)
	assert.NotNil(t, err)
}

func TestConnection_QueryContext7(t *testing.T) {
	t.Parallel()
	c := createConnectionFixture()

	e := c.Ping(context.Background())
	assert.Nil(t, e)

	driverRows, err := c.QueryContext(context.Background(), "StartQueryExecution_nil_error",
		[]driver.NamedValue{})
	assert.Nil(t, driverRows)
	assert.NotNil(t, err)

	dr, er := c.ExecContext(context.Background(), "StartQueryExecution_nil_error",
		[]driver.NamedValue{})
	assert.Nil(t, dr)
	assert.NotNil(t, er)

	driverRows, err = c.QueryContext(context.Background(), "StartQueryExecution_OK_GetQueryExecutionWithContext_QueryExecutionStateCancelled",
		[]driver.NamedValue{})
	assert.Nil(t, driverRows)
	assert.Equal(t, err, context.Canceled)

	driverRows, err = c.QueryContext(context.Background(), "StartQueryExecution_OK_GetQueryExecutionWithContext_QueryExecutionStateFailed",
		[]driver.NamedValue{})
	assert.Nil(t, driverRows)
	assert.Equal(t, err, ErrTestMockFailedByAthena)

	query := "SELECTExecContext_OK"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.Nil(t, er)
	assert.NotNil(t, dr)

	query = "SELECTQueryContext_OK"
	driverRows, err = c.QueryContext(context.Background(), query, []driver.NamedValue{})
	assert.Nil(t, err)
	assert.NotNil(t, driverRows)

	query = "SELECTQueryContext_OK"
	value := driver.NamedValue{Value: uint64(0)}
	driverRows, err = c.QueryContext(context.Background(), query, []driver.NamedValue{value})
	assert.NotNil(t, err)
	assert.Nil(t, driverRows)

	query = "SELECTQueryContext_?"
	value = driver.NamedValue{Value: "OK"}
	driverRows, err = c.QueryContext(context.Background(), query, []driver.NamedValue{value})
	assert.Nil(t, err)
	assert.NotNil(t, driverRows)

	query = "SELECTQueryContext_CANCEL_OK"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	driverRows, err = c.QueryContext(ctx, query, []driver.NamedValue{})
	assert.NotNil(t, err)
	assert.Nil(t, driverRows)

	query = "SELECTQueryContext_CANCEL_FAIL"
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	driverRows, err = c.QueryContext(ctx, query, []driver.NamedValue{})
	assert.NotNil(t, err)
	assert.Nil(t, driverRows)

	// After 3 seconds to get a timeout error
	query = "SELECTQueryContext_TIMEOUT"
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()
	driverRows, err = c.QueryContext(ctx, query, []driver.NamedValue{})
	assert.NotNil(t, err)
	assert.Nil(t, driverRows)

	query = randString(MAXQueryStringLength * 10)
	driverRows, err = c.QueryContext(context.Background(), query, []driver.NamedValue{})
	assert.Equal(t, err, ErrInvalidQuery)
	assert.Nil(t, driverRows)

	// Cancelled by AWS Athena
	query = "SELECTQueryContext_AWS_CANCEL"
	driverRows, err = c.QueryContext(context.Background(), query, []driver.NamedValue{})
	assert.NotNil(t, err)
	assert.Nil(t, driverRows)

	// failed by AWS Athena
	query = "SELECTQueryContext_AWS_FAIL"
	driverRows, err = c.QueryContext(context.Background(), query, []driver.NamedValue{})
	assert.NotNil(t, err)
	assert.Nil(t, driverRows)

	query = "SELECTQueryContext_CANCEL_OK"
	ctx, cancel = context.WithTimeout(context.Background(), PoolInterval*time.Second*2)
	defer cancel()
	driverRows, err = c.QueryContext(ctx, query, []driver.NamedValue{})
	assert.NotNil(t, err)
	assert.Nil(t, driverRows)

	query = "00000000-0000-0000-0000-000000000000"
	driverRows, err = c.QueryContext(context.Background(), query, []driver.NamedValue{})
	assert.Nil(t, err)
	assert.NotNil(t, driverRows)
}

func BenchmarkConnection_QueryContext(b *testing.B) {
	for i := 0; i < 10000; i++ {
		c := createConnectionFixture()
		assert.Nil(b, c.Ping(context.Background()))
	}
}

func createConnectionFixture() *Connection {
	rand.Seed(int64(time.Now().Nanosecond()))
	nm := newMockAthenaClient()
	c := &Connection{
		athenaAPI: nm,
		connector: NoopsSQLConnector(),
	}
	var s3bucket string = "s3://query-results-henry-wu-us-east-2/"
	wgTags := NewWGTags()
	wgTags.AddTag("Uber Author", "henry.wu")
	wgTags.AddTag("Uber Role", "Engineer")
	wg := NewDefaultWG("henry_wu", nil, wgTags)
	testConf := NewNoOpsConfig()
	_ = testConf.SetOutputBucket(s3bucket)
	_ = testConf.SetRegion(regions[rand.Int31n(int32(len(regions)))])
	testConf.SetUser("henry.wu")
	testConf.SetDB(randString(8)) // default
	testConf.SetWGRemoteCreationAllowed(true)
	nm.CreateWGStatus = true

	_ = testConf.SetWorkGroup(wg)
	c.connector.config = testConf
	return c
}

func TestMoneyWise(t *testing.T) {
	t.Parallel()
	c := &Connection{
		athenaAPI: newMockAthenaClient(),
		connector: NoopsSQLConnector(),
	}
	var s3bucket string = "s3://query-results-henry-wu-us-east-2/"

	wgTags := NewWGTags()
	wgTags.AddTag("Uber User", "henry.wu")
	wgTags.AddTag("Uber Asset", "abc.efg")
	wg := NewDefaultWG(DefaultWGName, nil, wgTags)
	testConf := NewNoOpsConfig()
	err := testConf.SetOutputBucket(s3bucket)
	assert.Nil(t, err)
	err = testConf.SetRegion("us-east-1")
	assert.Nil(t, err)
	testConf.SetUser("henry.wu")
	testConf.SetDB("default") // default

	_ = testConf.SetWorkGroup(wg)
	testConf.SetMoneyWise(true)
	c.connector.config = testConf
	query := "SELECTExecContext_OK"
	dr, er := c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.Nil(t, er)
	assert.NotNil(t, dr)

	query = "00000000-0000-0000-0000-000000000000"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.Nil(t, er)
	assert.NotNil(t, dr)

	query = "SELECTQueryContext_CANCEL_OK"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	dr2, err := c.QueryContext(ctx, query, []driver.NamedValue{})
	assert.NotNil(t, err)
	assert.Nil(t, dr2)

	query = "SELECTQueryContext_AWS_CANCEL"
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	dr2, err = c.QueryContext(ctx, query, []driver.NamedValue{})
	assert.NotNil(t, err)
	assert.Nil(t, dr2)
}

func TestConnection_CachedQuery(t *testing.T) {
	t.Parallel()
	c := &Connection{
		athenaAPI: newMockAthenaClient(),
		connector: NoopsSQLConnector(),
	}
	var s3bucket string = "s3://query-results-henry-wu-us-east-2/"
	testConf := NewNoOpsConfig()
	err := testConf.SetOutputBucket(s3bucket)
	assert.Nil(t, err)
	err = testConf.SetRegion("us-east-1")
	assert.Nil(t, err)
	testConf.SetUser("henry.wu")
	testConf.SetDB("default") // default
	testConf.SetMoneyWise(true)
	c.connector.config = testConf
	query := "00000000-0000-0000-0000-000000000000"
	dr, er := c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.Nil(t, er)
	assert.NotNil(t, dr)
}

func Test_PseudoCommand(t *testing.T) {
	t.Parallel()
	c := &Connection{
		athenaAPI: newMockAthenaClient(),
		connector: NoopsSQLConnector(),
	}
	var s3bucket string = "s3://query-results-henry-wu-us-east-2/"

	wgTags := NewWGTags()
	wgTags.AddTag("Uber User", "henry.wu")
	wgTags.AddTag("Uber Asset", "abc.efg")
	wg := NewDefaultWG(DefaultWGName, nil, wgTags)
	testConf := NewNoOpsConfig()
	err := testConf.SetOutputBucket(s3bucket)
	assert.Nil(t, err)
	err = testConf.SetRegion("us-east-1")
	assert.Nil(t, err)
	testConf.SetUser("henry.wu")
	testConf.SetDB("default") // default

	_ = testConf.SetWorkGroup(wg)
	testConf.SetMoneyWise(true)
	c.connector.config = testConf

	query := "pc:get_query_id"
	dr, er := c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.NotNil(t, er)
	assert.Nil(t, dr)

	query = "pc:get_query_id FAILED_AFTER_GETQID"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.Equal(t, er.Error(), "FAILED_AFTER_GETQID_FAILED")
	assert.Nil(t, dr)

	query = "pc:get_query_id FAILED_AFTER_GETQID2"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.Nil(t, er)
	assert.NotNil(t, dr)

	query = "pc:badcommand"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.NotNil(t, er)
	assert.Nil(t, dr)

	query = "pc:get_query_id SELECTQueryContext_CANCEL_OK"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.Nil(t, er)
	assert.NotNil(t, dr)

	query = "pc:stop_query_id c89088ab-595d-4ee6-a9ce-73b55aeb8954"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.Nil(t, er)
	assert.NotNil(t, dr)

	query = "pc:stop_query_id c89088ab-595d-4ee6-a9ce-73b55aeb8955"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.NotNil(t, er)
	assert.Nil(t, dr)

	query = "pc:stop_query_id"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.NotNil(t, er)
	assert.Nil(t, dr)

	query = "pc:stop_query_id 123"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.NotNil(t, er)
	assert.Nil(t, dr)

	query = "pc:get_query_id_status c89088ab-595d-4ee6-a9ce-73b55aeb8900"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.Nil(t, er)
	assert.NotNil(t, dr)

	query = "pc:get_query_id_status c89088ab-595d-4ee6-a9ce-73b55aeb8111"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.NotNil(t, er)
	assert.Nil(t, dr)

	query = "pc:get_query_id_status"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.NotNil(t, er)
	assert.Nil(t, dr)

	query = "pc:get_query_id_status 123"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.NotNil(t, er)
	assert.Nil(t, dr)

	query = "pc:get_driver_version"
	dr, er = c.ExecContext(context.Background(), query, []driver.NamedValue{})
	assert.Nil(t, er)
	assert.NotNil(t, dr)
}
