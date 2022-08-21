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
	"strings"

	drv "github.com/uber/athenadriver/go"
	"github.com/uber/athenadriver/lib/configfx"
	"github.com/uber/athenadriver/lib/queryfx"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(opts(), fx.Options(fx.NopLogger))
	ctx := context.Background()
	app.Start(ctx)
	defer app.Stop(ctx)
}

func opts() fx.Option {
	return fx.Options(
		configfx.Module,
		queryfx.Module,
		fx.Invoke(queryAthena),
	)
}

func queryAthena(qad queryfx.QueryAndDBConnection, mc configfx.AthenaDriverConfig) {
	for _, query := range qad.Query {
		query = strings.Trim(query, " \n\t")
		if query == "" {
			continue
		}
		rows, err := qad.DB.Query(query)
		if err != nil {
			println("ERROR: " + err.Error())
			if mc.OutputConfig.Fastfail {
				return
			}
			continue
		}
		defer rows.Close()
		if mc.OutputConfig.Rowonly {
			drv.PrettyPrintSQLRows(rows, mc.OutputConfig.Style, mc.OutputConfig.Render, mc.OutputConfig.Page)
		} else {
			drv.PrettyPrintSQLColsRows(rows, mc.OutputConfig.Style, mc.OutputConfig.Render, mc.OutputConfig.Page)
		}
	}
}
