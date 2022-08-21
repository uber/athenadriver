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
	"strings"
)

// Statement is to implement Go's database/sql Statement.
type Statement struct {
	connection *Connection
	closed     bool
	query      string
	numInput   int
}

// Close is to close an open statement.
func (s *Statement) Close() error {
	if s.connection == nil || s.closed {
		// driver.Stmt.Close can be called more than once, thus this function
		// has to be idempotent.
		// See also Issue #450 and golang/go#16019.
		return driver.ErrBadConn
	}
	s.query = ""
	s.closed = true
	s.numInput = 0
	return nil
}

// NumInput returns the number of prepared arguments.
// It may also return -1, if the driver doesn't know
// its number of placeholders. In that case, the sql package
// will not sanity check Exec or Query argument counts.
// -- From Go `sql/driver`
func (s *Statement) NumInput() int {
	if s.numInput == 0 {
		s.numInput = strings.Count(s.query, "?")
	}
	return s.numInput
}

// ColumnConverter is to return driver's DefaultParameterConverter.
func (s *Statement) ColumnConverter(idx int) driver.ValueConverter {
	return driver.DefaultParameterConverter
}

// Exec is to execute a prepared statement.
func (s *Statement) Exec(args []driver.Value) (driver.Result, error) {
	if s.closed {
		return nil, driver.ErrBadConn
	}
	r, e := s.connection.ExecContext(context.Background(), s.query,
		valueToNamedValue(args))
	s.closed = true
	return r, e
}

// Query is to query based on a prepared statement.
func (s *Statement) Query(args []driver.Value) (driver.Rows, error) {
	if s.closed {
		return nil, driver.ErrBadConn
	}
	r, e := s.connection.QueryContext(context.Background(), s.query,
		valueToNamedValue(args))
	s.closed = true
	return r, e
}
