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
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatement_NumInput(t *testing.T) {
	testConf := NewNoOpsConfig()
	connector := &SQLConnector{
		config: testConf,
		tracer: NewDefaultObservability(testConf),
	}

	conn, _ := connector.Connect(context.Background())
	st := Statement{
		connection: conn.(*Connection),
		query:      "abc=?",
	}
	dn := "123"
	d := []driver.Value{
		dn,
	}
	r, e := st.Exec(d)
	assert.NotNil(t, e)
	assert.Nil(t, r)
	assert.Equal(t, st.NumInput(), 1)
}

func TestStatement_Exec(t *testing.T) {
	testConf := NewNoOpsConfig()
	connector := &SQLConnector{
		config: testConf,
		tracer: NewDefaultObservability(testConf),
	}

	conn, _ := connector.Connect(context.Background())
	st := Statement{
		connection: conn.(*Connection),
		query:      "abc=?",
	}
	dn := "123"
	d := []driver.Value{
		dn,
	}
	_, e := st.Exec(d)
	assert.NotNil(t, e)
}

func TestStatement_Exec_After_Close(t *testing.T) {
	testConf := NewNoOpsConfig()
	connector := &SQLConnector{
		config: testConf,
		tracer: NewDefaultObservability(testConf),
	}

	conn, _ := connector.Connect(context.Background())
	st := Statement{
		connection: conn.(*Connection),
		query:      "abc=?",
	}
	dn := "123"
	d := []driver.Value{
		dn,
	}
	st.Close()
	_, err := st.Exec(d)
	assert.Equal(t, err, driver.ErrBadConn)
}

func TestStatement_Query(t *testing.T) {
	testConf := NewNoOpsConfig()
	connector := &SQLConnector{
		config: testConf,
		tracer: NewDefaultObservability(testConf),
	}

	conn, _ := connector.Connect(context.Background())
	st := Statement{
		connection: conn.(*Connection),
		query:      "abc=?",
	}
	dn := "123"
	d := []driver.Value{
		dn,
	}
	_, e := st.Query(d)
	assert.NotNil(t, e)
}

func TestStatement_Query_After_Close(t *testing.T) {
	testConf := NewNoOpsConfig()
	connector := &SQLConnector{
		config: testConf,
		tracer: NewDefaultObservability(testConf),
	}

	conn, _ := connector.Connect(context.Background())
	st := Statement{
		connection: conn.(*Connection),
		query:      "abc=?",
	}
	dn := "123"
	d := []driver.Value{
		dn,
	}
	st.Close()
	_, err := st.Query(d)
	assert.Equal(t, err, driver.ErrBadConn)
}

func TestStatement_ColumnConverter(t *testing.T) {
	testConf := NewNoOpsConfig()
	connector := &SQLConnector{
		config: testConf,
		tracer: NewDefaultObservability(testConf),
	}

	conn, _ := connector.Connect(context.Background())
	st := Statement{
		connection: conn.(*Connection),
		query:      "abc=?",
	}
	assert.NotNil(t, st.ColumnConverter(0))
}

func TestStatement_Close(t *testing.T) {
	testConf := NewNoOpsConfig()
	connector := &SQLConnector{
		config: testConf,
		tracer: NewDefaultObservability(testConf),
	}

	conn, _ := connector.Connect(context.Background())
	st := Statement{
		connection: conn.(*Connection),
		query:      "abc=?",
	}
	assert.Equal(t, st.NumInput(), 1)
	st.Close()
	assert.Equal(t, st.NumInput(), 0)
}

func TestStatement_Close_AfterConnectionClose(t *testing.T) {
	testConf := NewNoOpsConfig()
	connector := &SQLConnector{
		config: testConf,
		tracer: NewDefaultObservability(testConf),
	}

	conn, _ := connector.Connect(context.Background())
	st := Statement{
		connection: conn.(*Connection),
		query:      "abc=?",
	}
	conn.Close()
	st.connection = nil
	assert.Equal(t, st.NumInput(), 1)
	err := st.Close()
	assert.Equal(t, err, driver.ErrBadConn)
}

func TestStatement_Reuse(t *testing.T) {
	nm := newMockAthenaClient()
	testConf := NewNoOpsConfig()
	conn := &Connection{
		athenaAPI: nm,
		connector: NoopsSQLConnector(),
	}
	testConf.SetOutputBucket("s3://fake-query-results-arbitrary-bucket/")
	testConf.SetRegion("us-east-1")
	testConf.SetUser("test.user")
	testConf.SetDB("default")
	testConf.SetWGRemoteCreationAllowed(true)
	conn.connector.config = testConf
	st := Statement{
		connection: conn,
		query:      "SELECTQueryContext_?",
	}
	dn := "OK"
	d := []driver.Value{
		dn,
	}

	// First Query - should succeed with mock
	r1, e1 := st.Query(d)
	assert.Nil(t, e1) // First query succeeds
	assert.NotNil(t, r1)
	r1.Close()

	// Second Query - should also succeed (statement is reusable)
	r2, e2 := st.Query(d)
	assert.Nil(t, e2) // Second query also succeeds
	assert.NotNil(t, r2)
	r2.Close()

	// Third Query - should also succeed (statement remains reusable)
	r3, e3 := st.Query(d)
	assert.Nil(t, e3) // Third query also succeeds
	assert.NotNil(t, r3)
	r3.Close()

	// Also test Exec reuse
	r4, e4 := st.Exec(d)
	assert.Nil(t, e4) // Exec also succeeds
	assert.NotNil(t, r4)

	r5, e5 := st.Exec(d)
	assert.Nil(t, e5) // Second Exec also succeeds
	assert.NotNil(t, r5)
}

func TestStatement_Close_Idempotent(t *testing.T) {
	testConf := NewNoOpsConfig()
	connector := &SQLConnector{
		config: testConf,
		tracer: NewDefaultObservability(testConf),
	}

	conn, _ := connector.Connect(context.Background())
	st := Statement{
		connection: conn.(*Connection),
		query:      "abc=?",
	}

	// First Close
	err1 := st.Close()
	assert.Equal(t, err1, nil) // First close should succeed

	// Second Close - should be idempotent (return nil, not driver.ErrBadConn)
	err2 := st.Close()
	assert.Equal(t, err2, nil) // Should also be nil per driver contract
}
