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
