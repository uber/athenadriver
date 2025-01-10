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
	"database/sql/driver"
	credentials2 "github.com/aws/aws-sdk-go-v2/credentials"
	"os"
	"strconv"
	"time"

	"github.com/uber-go/tally"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/athena"
)

// SQLConnector is the connector for AWS Athena Driver.
type SQLConnector struct {
	config *Config
	tracer *DriverTracer
}

// NoopsSQLConnector is to create a noops SQLConnector.
func NoopsSQLConnector() *SQLConnector {
	noopsConfig := NewNoOpsConfig()
	return &SQLConnector{
		config: noopsConfig,
		tracer: NewDefaultObservability(noopsConfig),
	}
}

// AthenaClient is an interface to facilitate testing
type AthenaClient interface {
	CreateWorkGroup(context.Context, *athena.CreateWorkGroupInput, ...func(*athena.Options)) (*athena.CreateWorkGroupOutput, error)
	GetQueryExecution(context.Context, *athena.GetQueryExecutionInput, ...func(*athena.Options)) (*athena.GetQueryExecutionOutput, error)
	GetQueryResults(context.Context, *athena.GetQueryResultsInput, ...func(*athena.Options)) (*athena.GetQueryResultsOutput, error)
	GetWorkGroup(context.Context, *athena.GetWorkGroupInput, ...func(*athena.Options)) (*athena.GetWorkGroupOutput, error)
	StartQueryExecution(context.Context, *athena.StartQueryExecutionInput, ...func(options *athena.Options)) (*athena.StartQueryExecutionOutput, error)
	StopQueryExecution(context.Context, *athena.StopQueryExecutionInput, ...func(*athena.Options)) (*athena.StopQueryExecutionOutput, error)
}

// Driver is to construct a new SQLConnector.
func (c *SQLConnector) Driver() driver.Driver {
	return &SQLDriver{}
}

// Connect is to create an AWS session.
// The order to find auth information to create session is:
// 1. Manually set  AWS profile in Config by calling config.SetAWSProfile(profileName)
// 2. AWS_SDK_LOAD_CONFIG
// 3. Static Credentials
// Ref: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
func (c *SQLConnector) Connect(ctx context.Context) (driver.Conn, error) {
	now := time.Now()
	c.tracer = NewDefaultObservability(c.config)
	if metrics, ok := ctx.Value(MetricsKey).(tally.Scope); ok {
		c.tracer.SetScope(metrics)
	}
	if logger, ok := ctx.Value(LoggerKey).(*zap.Logger); ok {
		c.tracer.SetLogger(logger)
	}

	var awsCfg aws.Config
	var err error
	// respect AWS_SDK_LOAD_CONFIG and local ~/.aws/credentials, ~/.aws/config
	if ok, _ := strconv.ParseBool(os.Getenv("AWS_SDK_LOAD_CONFIG")); ok {
		if profile := c.config.GetAWSProfile(); profile != "" {
			awsCfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
			if err != nil {
				c.tracer.Scope().Counter(DriverName + ".failure.sqlconnector.newsession").Inc(1)
				return nil, err
			}
		}
	} else if c.config.GetAccessID() != "" {
		staticCredentials := credentials2.NewStaticCredentialsProvider(c.config.GetAccessID(),
			c.config.GetSecretAccessKey(),
			c.config.GetSessionToken())
		awsCfg = aws.Config{
			Region:      c.config.GetRegion(),
			Credentials: staticCredentials,
		}
	} else {
		awsCfg = aws.Config{
			Region: c.config.GetRegion(),
		}
	}

	athenaClient := athena.NewFromConfig(awsCfg)
	timeConnect := time.Since(now)
	conn := &Connection{
		athenaClient: athenaClient,
		connector:    c,
	}
	c.tracer.Scope().Timer(DriverName + ".connector.connect").Record(timeConnect)
	return conn, nil
}
