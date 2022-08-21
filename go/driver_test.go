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
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDriver(t *testing.T) {
	dsn := "s3://henry.wu%40uber.com:@query-results-henry-wu-us-east-2?db=default&" +
		"region=us-east-1&workgroup_config=%7B%0A++BytesScannedCutoffPerQuery%3A+1073741824%2C%0A++Enfo" +
		"rceWorkGroupConfiguration%3A+true%2C%0A++PublishCloudWatchMetricsEnabled%3A+true%2C%0A++Reques" +
		"terPaysEnabled%3A+false%0A%7D&workgroupName=henry_wu"
	pDB, err := sql.Open(DriverName, dsn)
	assert.Nil(t, err)
	assert.NotNil(t, pDB)

	pDB, err = sql.Open(DriverName+"x", "")
	assert.NotNil(t, err)
	assert.Nil(t, pDB)
}

func TestSQLDriver_Open(t *testing.T) {
	s := SQLDriver{
		conn: NoopsSQLConnector(),
	}
	testConf := NewNoOpsConfig()
	c, e := s.Open(testConf.Stringify())
	assert.Nil(t, e)
	assert.NotNil(t, c)

	s2 := SQLDriver{
		conn: NoopsSQLConnector(),
	}
	c, e = s2.Open("")
	assert.Nil(t, c)
	assert.NotNil(t, e)
}

func TestSQLDriver_OpenConnector(t *testing.T) {
	s := SQLDriver{
		conn: NoopsSQLConnector(),
	}
	testConf := NewNoOpsConfig()
	c, e := s.OpenConnector(testConf.Stringify())
	assert.Nil(t, e)
	assert.NotNil(t, c)
}
