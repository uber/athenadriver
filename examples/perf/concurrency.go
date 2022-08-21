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

package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	secret "github.com/uber/athenadriver/examples/constants"
	drv "github.com/uber/athenadriver/go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 1. Set AWS Credential in Driver Config.
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	conf, err := drv.NewDefaultConfig(
		"s3://qr-athena-query-result/",
		secret.Region,
		secret.AccessID,
		secret.SecretAccessKey)
	os.Setenv("AWS_REGION", "us-east-1")
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	numGoRoutine := 100
	wg.Add(numGoRoutine)
	for i := 0; i < numGoRoutine; i++ {
		// 2. Open Connection.
		db, _ := sql.Open(drv.DriverName, conf.Stringify())
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		//encoderCfg.TimeKey = ""
		atom := zap.NewAtomicLevel()
		logger := zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			atom,
		))
		atom.SetLevel(drv.DebugLevel)
		defer logger.Sync()
		go func(i int, conf *drv.Config) {
			defer wg.Done()
			// 3. Query cancellation after 2 seconds
			ctx := context.WithValue(context.Background(), drv.LoggerKey, logger)
			// 3. Query
			r, e := db.QueryContext(ctx, "SHOW FUNCTIONS")
			if e != nil {
				fmt.Errorf("[%v]%s\n", i, e.Error())
				return
			} else {
				print(i, ",")
			}
			defer r.Close()
			cnt := 0
			for r.Next() {
				cnt++
			}
		}(i, conf)

		go func(db *sql.DB, logger *zap.Logger) {
			for range time.Tick(2 * time.Second) {
				stats := db.Stats()
				logDBStats(stats, logger)
			}
		}(db, logger)
	}
	wg.Wait()
}

// logDBStats is to log DB statistics.
func logDBStats(stats sql.DBStats, logger *zap.Logger) {
	logger.Info("DBPoolStatus",
		zap.Int("MaxOpenConnections", stats.MaxOpenConnections),
		zap.Int("OpenConnections", stats.OpenConnections),
		zap.Int("Idle", stats.Idle),
		zap.Int("InUse", stats.InUse),
		zap.Int64("WaitCount", stats.WaitCount),
		zap.Duration("WaitDuration", stats.WaitDuration),
		zap.Int64("MaxIdleClosed", stats.MaxIdleClosed),
		zap.Int64("MaxLifetimeClosed", stats.MaxLifetimeClosed),
	)
}
