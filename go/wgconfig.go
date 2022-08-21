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

import "github.com/aws/aws-sdk-go/service/athena"

// WGConfig wraps WorkGroupConfiguration.
type WGConfig struct {
	wgConfig *athena.WorkGroupConfiguration
}

// GetDefaultWGConfig to create a default WorkGroupConfiguration.
func GetDefaultWGConfig() *athena.WorkGroupConfiguration {
	var bytesScannedCutoffPerQuery int64 = DefaultBytesScannedCutoffPerQuery
	var enforceWorkGroupConfiguration bool = true
	var publishCloudWatchMetricsEnabled bool = true
	var requesterPaysEnabled bool = false
	return &athena.WorkGroupConfiguration{
		BytesScannedCutoffPerQuery:      &bytesScannedCutoffPerQuery, // 1G by default
		EnforceWorkGroupConfiguration:   &enforceWorkGroupConfiguration,
		PublishCloudWatchMetricsEnabled: &publishCloudWatchMetricsEnabled,
		RequesterPaysEnabled:            &requesterPaysEnabled,
		ResultConfiguration:             nil,
	}
}

// NewWGConfig to create a WorkGroupConfiguration.
func NewWGConfig(bytesScannedCutoffPerQuery int64,
	enforceWorkGroupConfiguration bool,
	publishCloudWatchMetricsEnabled bool,
	requesterPaysEnabled bool,
	resultConfiguration *athena.ResultConfiguration) *athena.WorkGroupConfiguration {
	return &athena.WorkGroupConfiguration{
		BytesScannedCutoffPerQuery:      &bytesScannedCutoffPerQuery,
		EnforceWorkGroupConfiguration:   &enforceWorkGroupConfiguration,
		PublishCloudWatchMetricsEnabled: &publishCloudWatchMetricsEnabled,
		RequesterPaysEnabled:            &requesterPaysEnabled,
		ResultConfiguration:             resultConfiguration,
	}
}
