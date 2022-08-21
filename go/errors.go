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
	"errors"
	"fmt"
)

// Various errors the driver might return. Can change between driver versions.
var (
	ErrInvalidQuery                 = errors.New("query is not valid")
	ErrConfigInvalidConfig          = errors.New("driver config is invalid")
	ErrConfigOutputLocation         = errors.New("output location must starts with s3")
	ErrConfigRegion                 = errors.New("region is required")
	ErrConfigWGPointer              = errors.New("workgroup pointer is nil")
	ErrConfigAccessIDRequired       = errors.New("AWS access ID is required")
	ErrConfigAccessKeyRequired      = errors.New("AWS access Key is required")
	ErrQueryUnknownType             = errors.New("query parameter type is unknown")
	ErrQueryBufferOF                = errors.New("query buffer overflow")
	ErrQueryTimeout                 = errors.New("query timeout")
	ErrAthenaTransactionUnsupported = errors.New("Athena doesn't support transaction statements")
	ErrAthenaNilDatum               = errors.New("*athena.Datum must not be nil")
	ErrAthenaNilAPI                 = errors.New("athenaAPI must not be nil")
	ErrTestMockGeneric              = errors.New("some_mock_error_for_test")
	ErrTestMockFailedByAthena       = errors.New("the reason why Athena failed the query")
	ErrServiceLimitOverride         = fmt.Errorf("service limit override must be greater than %d", PoolInterval)
)
