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
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAthenaConfig(t *testing.T) {
	var s3bucket string = "s3://query-results-henry-wu-us-east-2/"

	wgTags := NewWGTags()
	wgTags.AddTag("Uber User", "henry.wu@uber.com")
	wgTags.AddTag("Uber Asset", "abc.efg")
	wg := NewDefaultWG("henry_wu", nil, wgTags)
	testConf := NewNoOpsConfig()
	err := testConf.SetOutputBucket(s3bucket)
	assert.Nil(t, err)
	err = testConf.SetRegion("us-east-1")
	assert.Nil(t, err)
	testConf.SetUser("henry.wu@uber.com")
	testConf.SetDB("default") // default

	err = testConf.SetWorkGroup(wg)
	assert.Nil(t, err)
	assert.Equal(t, testConf.GetUser(), "henry.wu@uber.com")
	assert.Equal(t, testConf.GetOutputBucket(), "s3://query-results-henry-wu-us-east-2/")
	expected := "s3://henry.wu%40uber.com:@query-results-henry-wu-us-east-2?WGRemoteCreation=true&db=default&missingAsEmptyString=true&region=us-east-1&tag=%7CUber+User%60henry.wu%40uber.com%7CUber+Asset%60abc.efg&workgroupConfig=%7B%0A++BytesScannedCutoffPerQuery%3A+1073741824%2C%0A++EnforceWorkGroupConfiguration%3A+true%2C%0A++PublishCloudWatchMetricsEnabled%3A+true%2C%0A++RequesterPaysEnabled%3A+false%0A%7D&workgroupName=henry_wu"
	actual := testConf.Stringify()
	assert.Equal(t, actual, expected)
	w := testConf.GetWorkgroup()
	assert.Equal(t, len(w.Tags.Get()), len(wgTags.Get()))

	x, err := NewConfig(expected)
	assert.Equal(t, x.GetOutputBucket(), s3bucket)
	assert.Nil(t, err)
}

func TestGetOutputBucket(t *testing.T) {
	var s3bucket string = "s3://query-results-henry-wu-us-east-2/local/"
	testConf := NewNoOpsConfig()
	err := testConf.SetOutputBucket(s3bucket)
	conf, _ := NewConfig(testConf.Stringify())
	assert.Nil(t, err)
	assert.Equal(t, testConf.GetOutputBucket(), "s3://query-results-henry-wu-us-east-2/local/")
	assert.Equal(t, conf.GetOutputBucket(), "s3://query-results-henry-wu-us-east-2/local/")
}

func TestAthenaConfigWrongS3Bucket(t *testing.T) {
	var s3bucket string = "file:///query-results-henry-wu-us-east-2/"
	testConf := NewNoOpsConfig()
	err := testConf.SetOutputBucket(s3bucket)
	assert.NotNil(t, err)
}

func TestConfig_SetOutputBucket(t *testing.T) {
	var s3bucket string = "s3://query-results-henry-wu-us-east-2"
	testConf := NewNoOpsConfig()
	err := testConf.SetOutputBucket(s3bucket)
	assert.Nil(t, err)
}

func TestAthenaConfigWrongRegion(t *testing.T) {
	testConf := NewNoOpsConfig()
	err := testConf.SetRegion("")
	assert.NotNil(t, err)
}

func TestAthenaConfigWrongWG(t *testing.T) {
	testConf := NewNoOpsConfig()
	err := testConf.SetWorkGroup(nil)
	assert.NotNil(t, err)

	wg := NewWG("wg", nil, nil)
	e := testConf.SetWorkGroup(wg)
	assert.Nil(t, e)
}

func TestAthenaConfigSafeString(t *testing.T) {
	var s3bucket string = "s3://query-results-henry-wu-us-east-2/"

	wg := NewDefaultWG("henry_wu", nil, nil)
	testConf := NewNoOpsConfig()
	err := testConf.SetOutputBucket(s3bucket)
	assert.Nil(t, err)
	err = testConf.SetRegion("us-east-1")
	assert.Nil(t, err)
	testConf.SetUser("henry.wu@uber.com")
	testConf.SetDB("default") // default
	err = testConf.SetWorkGroup(wg)
	assert.Nil(t, err)
	err = testConf.SetSecretAccessKey("thisisaKey")
	assert.Nil(t, err)
	err = testConf.SetAccessID("thisisanID")
	assert.Nil(t, err)
	testConf.SetSessionToken("thisisaToken")
	assert.Equal(t, testConf.GetUser(), "henry.wu@uber.com")
	assert.Equal(t, testConf.GetOutputBucket(), "s3://query-results-henry-wu-us-east-2/")
	expectedRawString := "s3://henry.wu%40uber.com:@query-results-henry-wu-us-east-2?WGRemoteCreation=true&accessID=thisisanID&db=default&missingAsEmptyString=true&region=us-east-1&secretAccessKey=thisisaKey&sessionToken=thisisaToken&tag=&workgroupConfig=%7B%0A++BytesScannedCutoffPerQuery%3A+1073741824%2C%0A++EnforceWorkGroupConfiguration%3A+true%2C%0A++PublishCloudWatchMetricsEnabled%3A+true%2C%0A++RequesterPaysEnabled%3A+false%0A%7D&workgroupName=henry_wu"
	expectedSafeString := "s3://henry.wu%40uber.com:@query-results-henry-wu-us-east-2?WGRemoteCreation=true&accessID=*&db=default&missingAsEmptyString=true&region=us-east-1&secretAccessKey=*&sessionToken=*&tag=&workgroupConfig=%7B%0A++BytesScannedCutoffPerQuery%3A+1073741824%2C%0A++EnforceWorkGroupConfiguration%3A+true%2C%0A++PublishCloudWatchMetricsEnabled%3A+true%2C%0A++RequesterPaysEnabled%3A+false%0A%7D&workgroupName=henry_wu"
	actualRaw := testConf.Stringify()
	actualSafe := testConf.SafeStringify()
	assert.Equal(t, expectedRawString, actualRaw)
	assert.Equal(t, expectedSafeString, actualSafe)

	x, err := NewConfig(expectedRawString)
	assert.Equal(t, x.GetOutputBucket(), s3bucket)
	assert.Nil(t, err)
}

func TestConfig_SetMaskedColumnValue(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetMaskedColumnValue("abc", "xxx")
	m, b := testConf.CheckColumnMasked("abc")
	assert.Equal(t, m, "xxx")
	assert.True(t, b)
	m, b = testConf.CheckColumnMasked("ABC")
	assert.NotEqual(t, m, "xxx")
	assert.False(t, b)
}

func TestConfig_SetMetrics(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetMetrics(true)
	assert.True(t, testConf.IsMetricsEnabled())
	testConf.SetMetrics(false)
	assert.False(t, testConf.IsMetricsEnabled())
}

func TestConfig_SetLogging(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetLogging(true)
	assert.True(t, testConf.IsLoggingEnabled())
	testConf.SetLogging(false)
	assert.False(t, testConf.IsLoggingEnabled())
}

func TestConfig_IsMissingAsEmptyString(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetMissingAsEmptyString(true)
	assert.True(t, testConf.IsMissingAsEmptyString())
	testConf.SetMissingAsEmptyString(false)
	assert.False(t, testConf.IsMissingAsEmptyString())
}

func TestConfig_IsMissingAsDefault(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetMissingAsDefault(true)
	assert.True(t, testConf.IsMissingAsDefault())
	testConf.SetMissingAsDefault(false)
	assert.False(t, testConf.IsMissingAsDefault())
}

func TestConfig_IsWGRemoteCreationAllowed(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetWGRemoteCreationAllowed(true)
	assert.True(t, testConf.IsWGRemoteCreationAllowed())
	testConf.SetWGRemoteCreationAllowed(false)
	assert.False(t, testConf.IsWGRemoteCreationAllowed())
}

func TestConfig_NewDefaultConfig(t *testing.T) {
	_, err := NewDefaultConfig("", "", "", "")
	assert.NotNil(t, err)
	_, err = NewDefaultConfig("file:///", "", "", "")
	assert.NotNil(t, err)
	_, err = NewDefaultConfig("s3:///abc", "", "", "")
	assert.NotNil(t, err)
	assert.NotNil(t, err)
	_, err = NewDefaultConfig("s3:///abc", "east", "", "")
	assert.NotNil(t, err)
	_, err = NewDefaultConfig("s3:///abc", "east", "as", "")
	assert.NotNil(t, err)
	_, err = NewDefaultConfig("s3:///abc", "east", "as", "ss")
	assert.Nil(t, err)
}

func TestConfig_NewConfig(t *testing.T) {
	x, err := NewConfig("\n")
	assert.NotNil(t, err)
	assert.Nil(t, x)
}

func TestConfig_GetWorkgroup(t *testing.T) {
	wg := NewDefaultWG("henry_wu", nil, nil)
	testConf := NewNoOpsConfig()
	err := testConf.SetWorkGroup(wg)
	assert.Nil(t, err)
	w := testConf.GetWorkgroup()
	assert.Nil(t, w.Tags)
}

func TestConfig_SetReadOnly(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetReadOnly(false)
	assert.False(t, testConf.IsReadOnly())
}

func TestConfig_GetDB(t *testing.T) {
	testConf := NewNoOpsConfig()
	assert.Equal(t, testConf.GetDB(), DefaultDBName)
	testConf.SetDB("")
	assert.Equal(t, testConf.GetDB(), DefaultDBName)
}

func TestConfig_GetRegion(t *testing.T) {
	testConf := NewNoOpsConfig()
	assert.Equal(t, testConf.GetRegion(), DefaultRegion)
	testConf = &Config{
		dsn:    *new(url.URL),
		values: url.Values{},
	}
	assert.Equal(t, testConf.GetRegion(), GetFromEnvVal(regionEnvKeys))
}

func TestConfig_GetAccessID(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetAccessID("abc")
	assert.Equal(t, testConf.GetAccessID(), "abc")
	testConf = &Config{
		dsn:    *new(url.URL),
		values: url.Values{},
	}
	assert.Equal(t, testConf.GetAccessID(), GetFromEnvVal(credAccessEnvKey))
}

func TestConfig_GetSecretAccessKey(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetSecretAccessKey("abc")
	assert.Equal(t, testConf.GetSecretAccessKey(), "abc")
	testConf = &Config{
		dsn:    *new(url.URL),
		values: url.Values{},
	}
	assert.Equal(t, testConf.GetSecretAccessKey(), GetFromEnvVal(credSecretEnvKey))
}

func TestConfig_GetSessionToken(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetSessionToken("abc")
	assert.Equal(t, testConf.GetSessionToken(), "abc")
	testConf = &Config{
		dsn:    *new(url.URL),
		values: url.Values{},
	}
	assert.Equal(t, testConf.GetSessionToken(), GetFromEnvVal(credSessionEnvKey))
}

func TestConfig_WGConfig(t *testing.T) {
	conf := NewWGConfig(10*DefaultBytesScannedCutoffPerQuery, true, true, false, nil)
	wg := NewDefaultWG("workgroup1", conf, nil)
	assert.Equal(t, *wg.Config.BytesScannedCutoffPerQuery, int64(DefaultBytesScannedCutoffPerQuery*10))
}

func TestConfig_SetMoneyWise(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetMoneyWise(false)
	assert.False(t, testConf.IsMoneyWise())
	testConf.SetMoneyWise(true)
	assert.True(t, testConf.IsMoneyWise())
}

func TestConfig_SetAWSProfile(t *testing.T) {
	testConf := NewNoOpsConfig()
	testConf.SetAWSProfile("development")
	assert.Equal(t, testConf.GetAWSProfile(), "development")
}

func TestConfig_SetServiceLimitOverride(t *testing.T) {
	var s3bucket string = "s3://query-results-henry-wu-us-east-2/"

	testConf := NewNoOpsConfig()
	_ = testConf.SetOutputBucket(s3bucket)
	serviceLimitOverride := NewServiceLimitOverride()
	ddlQueryTimeout := 1000 * 60 // 1000 minutes
	_ = serviceLimitOverride.SetDDLQueryTimeout(ddlQueryTimeout)
	testConf.SetServiceLimitOverride(*serviceLimitOverride)
	testServiceLimitOverride := testConf.GetServiceLimitOverride()
	assert.Equal(t, ddlQueryTimeout, testServiceLimitOverride.GetDDLQueryTimeout())

	expected := "s3://query-results-henry-wu-us-east-2?DDLQueryTimeout=60000&DMLQueryTimeout=0&WGRemoteCreation=true&db=default&missingAsEmptyString=true&region=us-east-1"
	assert.Equal(t, expected, testConf.Stringify())

	dmlQueryTimeout := 60 * 60 // 60 minutes
	_ = serviceLimitOverride.SetDMLQueryTimeout(dmlQueryTimeout)
	testConf.SetServiceLimitOverride(*serviceLimitOverride)
	testServiceLimitOverride = testConf.GetServiceLimitOverride()
	assert.Equal(t, ddlQueryTimeout, testServiceLimitOverride.GetDDLQueryTimeout())
	assert.Equal(t, dmlQueryTimeout, testServiceLimitOverride.GetDMLQueryTimeout())

	expected = "s3://query-results-henry-wu-us-east-2?DDLQueryTimeout=60000&DMLQueryTimeout=3600&WGRemoteCreation=true&db=default&missingAsEmptyString=true&region=us-east-1"
	assert.Equal(t, expected, testConf.Stringify())
}
