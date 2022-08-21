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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateTime_ScanTime(t *testing.T) {
	r, e := scanTime("01:02:03.456")
	assert.Nil(t, e)
	assert.True(t, r.Valid)
	assert.NotEqual(t, r.Time.String(), ZeroDateTimeString)
}

func TestDateTime_ScanTimeWithTimeZone(t *testing.T) {
	r, e := scanTime("01:02:03.456 America/Los_Angeles")
	assert.Nil(t, e)
	assert.True(t, r.Valid)
	assert.NotEqual(t, r.Time.String(), ZeroDateTimeString)

}

func TestDateTime_ScanTimeStamp(t *testing.T) {
	r, e := scanTime("2001-08-22 03:04:05.321")
	assert.Nil(t, e)
	assert.True(t, r.Valid)
	assert.NotEqual(t, r.Time.String(), ZeroDateTimeString)

}

func TestDateTime_ScanTimeStampWithTimeZone(t *testing.T) {
	r, e := scanTime("2001-08-22 03:04:05.321 America/Los_Angeles")
	assert.Nil(t, e)
	assert.True(t, r.Valid)
	assert.NotEqual(t, r.Time.String(), ZeroDateTimeString)
}

func TestDateTime_ScanTimeFail(t *testing.T) {
	r, e := scanTime("2001-08-22 03:04:05.321 PST")
	assert.NotNil(t, e)
	assert.False(t, r.Valid)
	assert.Equal(t, r.Time.String(), ZeroDateTimeString)

	r, e = scanTime("abc")
	assert.NotNil(t, e)
	assert.False(t, r.Valid)
	assert.Equal(t, r.Time.String(), ZeroDateTimeString)
}

func TestDateTime_ScanTimeFail_MonthOutOfRange(t *testing.T) {
	r, e := scanTime("2001-18-22 03:04:05.321 America/Los_Angeles")
	assert.NotNil(t, e)
	assert.False(t, r.Valid)
	assert.Equal(t, r.Time.String(), ZeroDateTimeString)
}

func TestDateTime_ParseAthenaTimeWithLocation(t *testing.T) {
	r, e := parseAthenaTimeWithLocation("abc")
	assert.NotNil(t, e)
	assert.False(t, r.Valid)
	assert.Equal(t, r.Time.String(), ZeroDateTimeString)

	r, e = parseAthenaTimeWithLocation("ab c")
	assert.NotNil(t, e)
	assert.False(t, r.Valid)
	assert.Equal(t, r.Time.String(), ZeroDateTimeString)
}
