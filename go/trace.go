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
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zap.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zap.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zap.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	// DPanicLevel, PanicLevel and FatalLevel are not allowed in this package
	// to avoid terminating the whole process
	ErrorLevel = zap.ErrorLevel
)

// DriverTracer is supported in athenadriver builtin.
type DriverTracer struct {
	logger *zap.Logger
	scope  tally.Scope
	config *Config
}

// NewObservability is to create an observability object.
func NewObservability(config *Config, logger *zap.Logger,
	scope tally.Scope) *DriverTracer {
	o := DriverTracer{
		logger: logger,
		scope:  scope,
		config: config,
	}
	return &o
}

// NewDefaultObservability is to create an observability object with logger
// and scope as default(noops object).
func NewDefaultObservability(config *Config) *DriverTracer {
	o := DriverTracer{
		logger: zap.NewNop(),
		scope:  tally.NoopScope,
		config: config,
	}
	return &o
}

// NewNoOpsObservability is for testing purpose.
func NewNoOpsObservability() *DriverTracer {
	o := DriverTracer{
		logger: zap.NewNop(),
		scope:  tally.NoopScope,
		config: NewNoOpsConfig(),
	}
	return &o
}

// Logger is a getter of logger.
func (c *DriverTracer) Logger() *zap.Logger {
	if !c.config.IsLoggingEnabled() {
		return zap.NewNop()
	}
	return c.logger
}

// SetLogger is a setter of logger.
func (c *DriverTracer) SetLogger(logger *zap.Logger) {
	c.logger = logger
}

// Scope is a getter of tally.Scope.
func (c *DriverTracer) Scope() tally.Scope {
	if !c.config.IsMetricsEnabled() {
		return tally.NoopScope
	}
	return c.scope
}

// SetScope is a setter of tally.Scope.
func (c *DriverTracer) SetScope(scope tally.Scope) {
	c.scope = scope
}

// Config is to get c.config
func (c *DriverTracer) Config() *Config {
	return c.config
}

// Log is to log with zap.logger with 4 logging levels.
// We threw away the panic and fatal level as we don't want to DB error terminates the whole process.
func (c *DriverTracer) Log(lvl zapcore.Level, msg string, fields ...zap.Field) {
	if !c.config.IsLoggingEnabled() {
		return
	}
	switch lvl {
	case DebugLevel:
		c.logger.Debug(msg, fields...)
	case WarnLevel:
		c.logger.Warn(msg, fields...)
	case InfoLevel:
		c.logger.Info(msg, fields...)
	case ErrorLevel, zap.DPanicLevel, zap.PanicLevel, zap.FatalLevel:
		c.logger.Error(msg, fields...)

	}
}
