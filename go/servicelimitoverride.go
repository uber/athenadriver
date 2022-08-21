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
	"fmt"
	"net/url"
	"strconv"
)

// ServiceLimitOverride allows users to override service limits, hardcoded in constants.go.
// This assumes the service limits have been raised in the AWS account.
// https://docs.aws.amazon.com/athena/latest/ug/service-limits.html
type ServiceLimitOverride struct {
	ddlQueryTimeout int
	dmlQueryTimeout int
}

// NewServiceLimitOverride is to create an empty ServiceLimitOverride.
// Values can be set using setters.
func NewServiceLimitOverride() *ServiceLimitOverride {
	return &ServiceLimitOverride{}
}

// SetDDLQueryTimeout is to set the DDLQueryTimeout override.
func (c *ServiceLimitOverride) SetDDLQueryTimeout(seconds int) error {
	if seconds < PoolInterval {
		return ErrServiceLimitOverride
	}
	c.ddlQueryTimeout = seconds
	return nil
}

// GetDDLQueryTimeout is to get the DDLQueryTimeout override.
func (c *ServiceLimitOverride) GetDDLQueryTimeout() int {
	return c.ddlQueryTimeout
}

// SetDMLQueryTimeout is to set the DMLQueryTimeout override.
func (c *ServiceLimitOverride) SetDMLQueryTimeout(seconds int) error {
	if seconds < PoolInterval {
		return ErrServiceLimitOverride
	}
	c.dmlQueryTimeout = seconds
	return nil
}

// GetDMLQueryTimeout is to get the DMLQueryTimeout override.
func (c *ServiceLimitOverride) GetDMLQueryTimeout() int {
	return c.dmlQueryTimeout
}

// GetAsStringMap is to get the ServiceLimitOverride as a map of strings
// and aids in setting url.Values in Config
func (c *ServiceLimitOverride) GetAsStringMap() map[string]string {
	res := map[string]string{}
	res["DDLQueryTimeout"] = fmt.Sprintf("%d", c.ddlQueryTimeout)
	res["DMLQueryTimeout"] = fmt.Sprintf("%d", c.dmlQueryTimeout)
	return res
}

// SetFromValues is to set ServiceLimitOverride properties from a url.Values
// which might be a list of override and other ignored values from a dsn
func (c *ServiceLimitOverride) SetFromValues(kvp url.Values) {
	ddlQueryTimeout, _ := strconv.Atoi(kvp.Get("DDLQueryTimeout"))
	_ = c.SetDDLQueryTimeout(ddlQueryTimeout)
	dmlQueryTimeout, _ := strconv.Atoi(kvp.Get("DMLQueryTimeout"))
	_ = c.SetDMLQueryTimeout(dmlQueryTimeout)
}
