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
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

func TestObservability_Config(t *testing.T) {
	obs := NewNoOpsObservability()
	assert.Equal(t, obs.Config(), NewNoOpsConfig())
}

func TestObservability_Scope(t *testing.T) {
	obs := NewNoOpsObservability()
	assert.Equal(t, obs.Scope(), tally.NoopScope)

	config := NewNoOpsConfig()
	config.SetMetrics(true)
	obs = NewDefaultObservability(config)
	assert.Equal(t, obs.Scope(), tally.NoopScope)
}

func TestObservability_Logger(t *testing.T) {
	obs := NewNoOpsObservability()
	assert.Equal(t, obs.Logger(), zap.NewNop())

	config := NewNoOpsConfig()
	config.SetLogging(false)
	obs = NewDefaultObservability(config)
	assert.Equal(t, obs.Logger(), zap.NewNop())
}

func TestObservability_Log(t *testing.T) {
	config := NewNoOpsConfig()
	config.SetLogging(false)
	obs := NewDefaultObservability(config)
	obs.Log(-1, "")
	config.SetLogging(true)
	obs = NewDefaultObservability(config)
	obs.Log(-1, "")
	obs.Log(ErrorLevel, "")
	obs.Log(WarnLevel, "")
	obs.Log(InfoLevel, "")
	obs.Log(DebugLevel, "")
}

func TestObservability_SetScope(t *testing.T) {
	obs := NewNoOpsObservability()
	obs.SetScope(tally.NoopScope)
	assert.Equal(t, obs.Scope(), tally.NoopScope)
}

func TestObservability_SetLogger(t *testing.T) {
	obs := NewNoOpsObservability()
	obs.SetLogger(nil)
	assert.Nil(t, obs.Logger())
}

func TestObservability_NewObservability(t *testing.T) {
	obs := NewObservability(NewNoOpsConfig(), zap.NewNop(), tally.NoopScope)
	assert.NotNil(t, obs.Logger())
	assert.Equal(t, obs.Logger(), zap.NewNop())
}
