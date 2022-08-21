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
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/stretchr/testify/assert"
)

// variadicToSlice, https://blog.learngoprogramming.com/golang-variadic-funcs-how-to-patterns-369408f19085
// https://stackoverflow.com/questions/23723955/how-can-i-pass-a-slice-as-a-variadic-input
// If f is variadic with final parameter type ...T, then within the function
// the argument is equivalent to a parameter of type []T. At each call of f,
// the argument passed to the final parameter is a new slice of type []T whose
// successive elements are the actual arguments,
// which all must be assignable to the type T.
func variadicToSlice(dest ...driver.Value) []driver.Value {
	return dest
}

func TestOnePageSuccess(t *testing.T) {
	testConf := NewNoOpsConfig()
	tests := []struct {
		desc                string
		queryID             string
		expectedResultsSize int
		expectedError       error
	}{
		{
			desc:                "show query, header, 5 row, no error",
			queryID:             "show",
			expectedResultsSize: 5,
			expectedError:       nil,
		},
	}
	for _, test := range tests {
		r, _ := NewRows(context.Background(), newMockAthenaClient(),
			test.queryID, testConf, NewDefaultObservability(testConf))

		var testArray, firstName, lastName string
		var active bool
		var uid int
		var registerDate, registerTS time.Time
		cnt := 0
		var err error = nil
		for {
			err = r.Next(variadicToSlice(&testArray, &active, &firstName, &lastName,
				&uid, &registerDate, &registerTS))
			if err != nil {
				if err != io.EOF {
					assert.Equal(t, test.expectedError, err)
				}
				break
			}
			cnt++
		}
		assert.Equal(t, test.expectedResultsSize, cnt)
		if err != io.EOF {
			assert.Equal(t, test.expectedError, err)
		}
		r.Close()
	}
}

func TestNextFailure(t *testing.T) {
	testConf := NewNoOpsConfig()
	tests := []struct {
		desc                string
		queryID             string
		expectedResultsSize int
		expectedError       error
	}{
		{
			desc:                "failed during calling next",
			queryID:             "RowsNextFailed",
			expectedResultsSize: 4,
			expectedError:       ErrTestMockGeneric,
		},
	}
	for _, test := range tests {
		r, _ := NewRows(context.Background(), newMockAthenaClient(),
			test.queryID,
			testConf, NewDefaultObservability(testConf))

		var testArray, firstName, lastName string
		var active bool
		var uid int
		var registerDate, registerTS time.Time
		cnt := 0
		var err error = nil
		for {
			err = r.Next(variadicToSlice(&testArray, &active, &firstName, &lastName,
				&uid, &registerDate, &registerTS))
			if err != nil {
				if err != io.EOF {
					assert.Equal(t, test.expectedError, err)
				}
				break
			}
			cnt++
		}
		assert.Equal(t, test.expectedResultsSize, cnt)
		if err != io.EOF {
			assert.Equal(t, test.expectedError, err)
		}
	}
}

func TestMultiplePages(t *testing.T) {
	testConf := NewNoOpsConfig()
	tests := []struct {
		desc                string
		queryID             string
		expectedResultsSize int
		expectedError       error
	}{
		{
			desc:                "select query, header, multiple pages",
			queryID:             "SELECT_OK",
			expectedResultsSize: 35,
			expectedError:       nil,
		},
	}
	var r *Rows
	for _, test := range tests {
		r, _ = NewRows(context.Background(), newMockAthenaClient(),
			test.queryID,
			testConf, NewDefaultObservability(testConf))

		var testArray, firstName, lastName string
		var active bool
		var uid int
		var registerDate, registerTS time.Time
		cnt := 0
		var err error = nil
		for {
			err = r.Next(variadicToSlice(&testArray, &active, &firstName, &lastName,
				&uid, &registerDate, &registerTS))
			if err != nil {
				if err != io.EOF {
					assert.Equal(t, test.expectedError, err)
				}
				break
			}
			cnt++
		}
		assert.Equal(t, test.expectedResultsSize, cnt)
		if err != io.EOF {
			assert.Equal(t, test.expectedError, err)
		}
	}
	var dest []driver.Value = make([]driver.Value, 8)
	assert.Equal(t, r.Next(dest), io.EOF)
}

func TestRows_Columns(t *testing.T) {
	testConf := NewNoOpsConfig()
	cs := createTestColumns()
	tests := []struct {
		desc                string
		queryID             string
		expectedResultsSize int
		expectedError       error
	}{
		{
			desc:                "select query, header, multiple pages",
			queryID:             "SELECT_OK",
			expectedResultsSize: 35,
			expectedError:       nil,
		},
	}
	for _, test := range tests {
		r, _ := NewRows(context.Background(), newMockAthenaClient(),
			test.queryID,
			testConf, NewDefaultObservability(testConf))
		assert.Equal(t, len(r.Columns()), len(cs))
	}
}

func TestRows_ColumnTypeDatabaseTypeName(t *testing.T) {
	testConf := NewNoOpsConfig()
	cs := createTestColumns()
	tests := []struct {
		desc                string
		queryID             string
		expectedResultsSize int
		expectedError       error
	}{
		{
			desc:                "select query, header, multiple pages",
			queryID:             "SELECT_OK",
			expectedResultsSize: 35,
			expectedError:       nil,
		},
	}
	for _, test := range tests {
		r, _ := NewRows(context.Background(), newMockAthenaClient(),
			test.queryID,
			testConf, NewDefaultObservability(testConf))
		for i, v := range cs {
			assert.Equal(t, r.ColumnTypeDatabaseTypeName(i), *v.Type)

		}

	}
}

func TestRows_GetDefaultValueForColumnType(t *testing.T) {
	testConf := NewNoOpsConfig()
	tests := []struct {
		desc                string
		queryID             string
		expectedResultsSize int
		expectedError       error
	}{
		{
			desc:                "select query, header, multiple pages",
			queryID:             "SELECT_OK",
			expectedResultsSize: 35,
			expectedError:       nil,
		},
	}
	for _, test := range tests {
		r, _ := NewRows(context.Background(), newMockAthenaClient(),
			test.queryID,
			testConf, NewDefaultObservability(testConf))
		for _, v := range []string{"tinyint", "smallint", "integer", "bigint"} {
			assert.Equal(t, r.getDefaultValueForColumnType(v), 0)
		}
		for _, v := range []string{"json", "char", "varchar", "varbinary", "row", "string", "binary",
			"struct", "interval year to month", "interval day to second", "decimal",
			"ipaddress", "array", "map", "unknown"} {
			assert.Equal(t, r.getDefaultValueForColumnType(v), "")
		}
		for _, v := range []string{"float", "double", "real"} {
			assert.Equal(t, r.getDefaultValueForColumnType(v), 0.0)
		}
		for _, v := range []string{"date", "time", "time with time zone", "timestamp", "timestamp with time zone"} {
			assert.Equal(t, r.getDefaultValueForColumnType(v), time.Time{})
		}
		assert.Equal(t, r.getDefaultValueForColumnType("boolean"), false)
		assert.Equal(t, r.getDefaultValueForColumnType("XXX"), "")
	}
}

func TestRows_AthenaTypeToGoType(t *testing.T) {
	testConf := NewNoOpsConfig()
	r, _ := NewRows(context.Background(), newMockAthenaClient(),
		"SELECT_OK", testConf, NewDefaultObservability(testConf))
	c := newColumnInfo("a", "tinyint")
	// tinyint
	rv := "1"
	g, e := r.athenaTypeToGoType(c, &rv, testConf)
	assert.Nil(t, e)
	assert.Equal(t, int8(1), g)

	rv = "x"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.NotNil(t, e)
	assert.Nil(t, g)

	// smallint
	c = newColumnInfo("a", "smallint")
	rv = "1"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.Nil(t, e)
	assert.Equal(t, int16(1), g)

	rv = "x"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.NotNil(t, e)
	assert.Nil(t, g)

	// int
	c = newColumnInfo("a", "integer")
	rv = "1"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.Nil(t, e)
	assert.Equal(t, int32(1), g)

	rv = "x"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.NotNil(t, e)
	assert.Nil(t, g)

	// bigint
	c = newColumnInfo("a", "bigint")
	rv = "1"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.Nil(t, e)
	assert.Equal(t, int64(1), g)

	rv = "x"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.NotNil(t, e)
	assert.Nil(t, g)

	// float
	c = newColumnInfo("a", "float")
	rv = "1.0"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.Nil(t, e)
	assert.Equal(t, float32(1.0), g)

	rv = "x"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.NotNil(t, e)
	assert.Nil(t, g)

	// real
	c = newColumnInfo("a", "real")
	rv = "1.0"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.Nil(t, e)
	assert.Equal(t, float32(1.0), g)

	rv = "x"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.NotNil(t, e)
	assert.Nil(t, g)

	// double
	c = newColumnInfo("a", "double")
	rv = "1.0"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.Nil(t, e)
	assert.Equal(t, float64(1.0), g)

	rv = "x"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.NotNil(t, e)
	assert.Nil(t, g)

	// string-like
	for _, s := range []string{"json", "char", "varchar", "varbinary", "row",
		"string", "binary",
		"struct", "interval year to month", "interval day to second", "decimal",
		"ipaddress", "array", "map", "unknown"} {
		c = newColumnInfo("a", s)
		rv = "012"
		g, e = r.athenaTypeToGoType(c, &rv, testConf)
		assert.Nil(t, e)
		assert.Equal(t, "012", g)
	}

	// boolean
	for _, s := range []string{"boolean"} {
		c = newColumnInfo("a", s)
		rv = "true"
		g, e = r.athenaTypeToGoType(c, &rv, testConf)
		assert.Nil(t, e)
		assert.Equal(t, true, g)

		rv = "false"
		g, e = r.athenaTypeToGoType(c, &rv, testConf)
		assert.Nil(t, e)
		assert.Equal(t, false, g)

		rv = "x"
		g, e = r.athenaTypeToGoType(c, &rv, testConf)
		assert.NotNil(t, e)
		assert.Nil(t, g)
	}

	// date and time
	now := time.Now()
	for _, s := range []string{"date", "time", "time with time zone",
		"timestamp", "timestamp with time zone"} {
		c = newColumnInfo("a", s)
		rv = "2020-01-20"
		g, e = r.athenaTypeToGoType(c, &rv, testConf)
		assert.Nil(t, e)
		assert.Equal(t, reflect.TypeOf(now), reflect.TypeOf(g))

		rv = "x"
		g, e = r.athenaTypeToGoType(c, &rv, testConf)
		assert.NotNil(t, e)
		assert.Nil(t, g)
	}

	c = newColumnInfo("a", "some_weird_type")
	rv = "123"
	g, e = r.athenaTypeToGoType(c, &rv, testConf)
	assert.NotNil(t, e)
	assert.Nil(t, g)

	// missing data - rawValue is nil
	c = newColumnInfo("a", "integer")
	g, e = r.athenaTypeToGoType(c, nil, testConf)
	assert.Nil(t, e)
	assert.Equal(t, g, "")

	testConf.SetMissingAsEmptyString(false)
	testConf.SetMissingAsDefault(true)
	g, e = r.athenaTypeToGoType(c, nil, testConf)
	assert.Nil(t, e)
	assert.Equal(t, g, 0)

	testConf.SetMissingAsEmptyString(false)
	testConf.SetMissingAsDefault(false)
	g, e = r.athenaTypeToGoType(c, nil, testConf)
	assert.NotNil(t, e)
	assert.Nil(t, g)

	// masked column
	testConf.SetMaskedColumnValue("a", "xxx")
	g, e = r.athenaTypeToGoType(c, nil, testConf)
	assert.Nil(t, e)
	assert.Equal(t, g, "xxx")
}

func TestRows_ColumnTypeDatabaseTypeName2(t *testing.T) {
	testConf := NewNoOpsConfig()
	r, _ := NewRows(context.Background(), newMockAthenaClient(),
		"SELECT_OK", testConf, NewDefaultObservability(testConf))
	c := newColumnInfo("a", nil)
	getQueryResultsOutput := &athena.GetQueryResultsOutput{
		ResultSet: &athena.ResultSet{
			ResultSetMetadata: &athena.ResultSetMetadata{
				ColumnInfo: []*athena.ColumnInfo{
					c,
				},
			},
		},
	}
	r.ResultOutput = getQueryResultsOutput
	assert.Equal(t, r.ColumnTypeDatabaseTypeName(0), "")
}

func TestRows_NewRows(t *testing.T) {
	testConf := NewNoOpsConfig()
	r, e := NewRows(context.Background(), newMockAthenaClient(),
		"1coloumn0row",
		testConf, NewDefaultObservability(testConf))
	assert.Nil(t, e)
	assert.NotNil(t, r)

	r, e = NewRows(context.Background(), newMockAthenaClient(),
		"1coloumn0row_valid",
		testConf, NewDefaultObservability(testConf))
	assert.Nil(t, e)
	assert.Equal(t, *r.ResultOutput.ResultSet.Rows[0].Data[0].VarCharValue,
		"1024")

	r, e = NewRows(context.Background(), newMockAthenaClient(),
		"column_more_than_row_fields",
		testConf, NewDefaultObservability(testConf))
	assert.Nil(t, e)
	assert.NotNil(t, r)

	r, e = NewRows(context.Background(), newMockAthenaClient(),
		"row_fields_more_than_column",
		testConf, NewDefaultObservability(testConf))
	assert.Nil(t, e)
	assert.NotNil(t, r)

	r, e = NewRows(context.Background(), newMockAthenaClient(),
		"GetQueryResultsWithContext_return_error",
		testConf, NewDefaultObservability(testConf))
	assert.NotNil(t, e)
	assert.Nil(t, r)

	// rawValue is nil
	r, e = NewRows(context.Background(), newMockAthenaClient(),
		"missing_data_resp",
		testConf, NewDefaultObservability(testConf))
	assert.Nil(t, e)
	assert.NotNil(t, r)
	var dest []driver.Value = make([]driver.Value, 8)
	e = r.Next(dest)
	assert.Equal(t, e, nil)

	// raise error for missing value
	testConf.SetMissingAsEmptyString(false)
	testConf.SetMissingAsDefault(false)
	r, e = NewRows(context.Background(), newMockAthenaClient(),
		"missing_data_resp",
		testConf, NewDefaultObservability(testConf))
	assert.Nil(t, e)
	assert.NotNil(t, r)
	e = r.Next(dest)
	assert.Equal(t, e.Error(), "Missing data at column c1")

	r, e = NewRows(context.Background(), newMockAthenaClient(),
		"missing_data_resp2",
		testConf, NewDefaultObservability(testConf))
	assert.Nil(t, e)
	assert.NotNil(t, r)
	e = r.Next(dest)
	assert.NotEqual(t, e, io.EOF)

	// error when row.Next()
	r, e = NewRows(context.Background(), newMockAthenaClient(),
		"SELECT_GetQueryResults_ERR",
		testConf, NewDefaultObservability(testConf))
	assert.Nil(t, e)
	assert.NotNil(t, r)
	for {
		e = r.Next(dest)
		if e != nil {
			assert.Equal(t, e, ErrTestMockGeneric)
			break
		}
	}

	// missing row in page
	r, e = NewRows(context.Background(), newMockAthenaClient(),
		"SELECT_EMPTY_ROW_IN_PAGE",
		testConf, NewDefaultObservability(testConf))
	assert.Nil(t, e)
	assert.NotNil(t, r)
	for {
		e = r.Next(dest)
		if e != nil {
			assert.Equal(t, e, io.EOF)
			break
		}
	}

	// close in the loop
	r, e = NewRows(context.Background(), newMockAthenaClient(),
		"SELECT_GetQueryResults_ERR",
		testConf, NewDefaultObservability(testConf))
	assert.Nil(t, e)
	assert.NotNil(t, r)
	cnt := 0
	for {
		e = r.Next(dest)
		if e != nil {
			assert.Equal(t, e, io.EOF)
			break
		}
		if cnt == 7 {
			r.Close()
		}
		cnt++
	}

}
