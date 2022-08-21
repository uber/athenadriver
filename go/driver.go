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
)

// SQLDriver is an implementation of sql/driver interface for AWS Athena.
// https://vyskocilm.github.io/blog/implement-sql-database-driver-in-100-lines-of-go/
// https://golang.org/pkg/database/sql/driver/#Driver
type SQLDriver struct {
	conn *SQLConnector
}

func init() {
	sql.Register(DriverName, &SQLDriver{})
}

// Open returns a new connection to AWS Athena.
// The dsn is a string in a driver-specific format.
// the sql package maintains a pool of idle connections for efficient re-use.
// The returned connection is only used by one goroutine at a time.
func (d *SQLDriver) Open(dsn string) (driver.Conn, error) {
	config, err := NewConfig(dsn)
	if err != nil {
		return nil, err
	}
	c := &SQLConnector{
		config: config,
	}
	return c.Connect(context.Background())
}

// OpenConnector will be called upon query execution.
// If a Driver implements DriverContext.OpenConnector, then sql.DB will call
// OpenConnector to obtain a Connector and then invoke
// that Connector's Conn method to obtain each needed connection,
// instead of invoking the Driver's Open method for each connection.
// The two-step sequence allows drivers to parse the name just once
// and also provides access to per-Conn contexts.
func (d *SQLDriver) OpenConnector(dsn string) (driver.Connector, error) {
	config, err := NewConfig(dsn)
	d.conn = &SQLConnector{
		config: config,
	}
	return d.conn, err
}
