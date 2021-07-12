package athenadriver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Tests for ServiceLimitOverride.
func TestNewServiceLimitOverride(t *testing.T) {
	testConf := NewServiceLimitOverride()
	assert.Zero(t, testConf.GetDDLQueryTimeout())
	assert.Zero(t, testConf.GetDMLQueryTimeout())

	ddlQueryTimeout := 30 * 60 // seconds
	dmlQueryTimeout := 60 * 60 // seconds
	testConf.SetDDLQueryTimeout(ddlQueryTimeout)
	assert.Equal(t, ddlQueryTimeout, testConf.GetDDLQueryTimeout()) // seconds

	testConf.SetDMLQueryTimeout(dmlQueryTimeout)
	assert.Equal(t, dmlQueryTimeout, testConf.GetDMLQueryTimeout()) // seconds

	ddlQueryTimeout = 0
	dmlQueryTimeout = 0
	err := testConf.SetDDLQueryTimeout(ddlQueryTimeout)
	assert.NotNil(t, err)

	err = testConf.SetDMLQueryTimeout(dmlQueryTimeout)
	assert.NotNil(t, err)

	ddlQueryTimeout = -1
	dmlQueryTimeout = -1
	err = testConf.SetDDLQueryTimeout(ddlQueryTimeout)
	assert.NotNil(t, err)

	err = testConf.SetDMLQueryTimeout(dmlQueryTimeout)
	assert.NotNil(t, err)
}
