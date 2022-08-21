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

package configfx

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	secret "github.com/uber/athenadriver/examples/constants"
	drv "github.com/uber/athenadriver/go"
	"go.uber.org/config"
	"go.uber.org/fx"
)

// Module is to provide dependency of Configuration to main app
var Module = fx.Provide(new)

// Params defines the dependencies or inputs
type Params struct {
	fx.In

	LC fx.Lifecycle
}

// ReaderOutputConfig is to represent the output section of configuration file
type ReaderOutputConfig struct {
	// Render is for the output format
	Render string `yaml:"render"`
	// Page is for the pagination
	Page int `yaml:"pagesize"`
	// Style is output style
	Style string `yaml:"style"`
	// Rowonly is for displaying header or not
	Rowonly bool `yaml:"rowonly"`
	// Moneywise is for displaying spending or not
	Moneywise bool `yaml:"moneywise"`
	// Fastfail is for multiple queries
	Fastfail bool `yaml:"fastfail"`
}

// ReaderInputConfig is to represent the input section of configuration file
type ReaderInputConfig struct {
	// Bucket is the output bucket
	Bucket string `yaml:"bucket"`
	// Region is AWS region
	Region string `yaml:"region"`
	// Database is the name of the DB
	Database string `yaml:"database"`
	// Admin is for write mode
	Admin bool `yaml:"admin"`
}

// AthenaDriverConfig is Athena Driver Configuration
type AthenaDriverConfig struct {
	// OutputConfig is for the output section of the config
	OutputConfig ReaderOutputConfig
	// InputConfig is for the input section of the config
	InputConfig ReaderInputConfig
	// QueryString is the query string
	QueryString []string
	// DrvConfig is the datastructure of Driver Config
	DrvConfig *drv.Config
}

// Result defines output
type Result struct {
	fx.Out

	// MyConfig is the current AthenaDriver Config
	MyConfig AthenaDriverConfig
}

func init() {
	setUpFlagUsage(context.Background())
}

func new(p Params) (Result, error) {
	var mc = AthenaDriverConfig{
		QueryString: make([]string, 0),
	}
	var (
		provider *config.YAML
		err      error
	)

	p.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			os.Unsetenv("AWS_SDK_LOAD_CONFIG")
			return nil
		},
	})

	var bucket = flag.String("b", secret.OutputBucket, "Athena resultset output bucket")
	var database = flag.String("d", "default", "The database you want to query")
	var query = flag.String("q", "select 1", "The SQL query string or a file containing SQL string")
	var rowOnly = flag.Bool("r", false, "Display rows only, don't show the first row as columninfo")
	var moneyWise = flag.Bool("m", false, "Enable moneywise mode to display the query cost as the first line of the output")
	var versionFlag = flag.Bool("v", false, "Print the current version and exit")
	var admin = flag.Bool("a", false, "Enable admin mode, so database write(create/drop) is allowed at athenadriver level")
	var style = flag.String("y", "default", "Output rendering style")
	var format = flag.String("o", "csv", "Output format(options: table, markdown, csv, html)")
	var fastFail = flag.Bool("f", true, "fast fail when where are multiple queries")

	flag.Parse()
	switch {
	case *versionFlag:
		println("Current build version: v" + drv.DriverVersion)
		os.Exit(0)
		return Result{}, fmt.Errorf("no")
	}

	// How to install a config file from a library
	if _, err = os.Stat(homeDir() + "/athenareader.config"); err == nil {
		provider, err = config.NewYAML(config.File(homeDir() + "/athenareader.config"))
	} else if _, err = os.Stat("athenareader.config"); err == nil {
		provider, err = config.NewYAML(config.File("athenareader.config"))
	} else {
		goPath := os.Getenv("GOPATH")
		if goPath == "" {
			goPath = homeDir() + "/go"
			if _, err = os.Stat(goPath); err != nil {
				d, _ := os.Getwd()
				println("could not find athenareader.config in home directory or current directory " + d)
				os.Exit(1)
			}
		}
		path := goPath + "/src/github.com/uber/athenadriver/athenareader/athenareader.config"
		if _, err = os.Stat(path); err == nil {
			copyFile(path, homeDir()+"/athenareader.config")
			provider, err = config.NewYAML(config.File(path))
		} else {
			err = downloadFile(homeDir()+"/athenareader.config",
				"https://raw.githubusercontent.com/uber/athenadriver/master/athenareader/athenareader.config")
			if err != nil {
				d, _ := os.Getwd()
				println("could not find athenareader.config in home directory or current directory " + d)
				os.Exit(1)
			} else {
				provider, err = config.NewYAML(config.File(homeDir() + "/athenareader.config"))
			}
		}
	}

	if err != nil {
		return Result{}, err
	}

	provider.Get("athenareader.output").Populate(&mc.OutputConfig)
	provider.Get("athenareader.input").Populate(&mc.InputConfig)

	filePath := expand(*query)
	if _, err := os.Stat(filePath); err == nil {
		b, err := ioutil.ReadFile(filePath)
		if err == nil {
			mc.QueryString = strings.Split(string(b), "\n\n") // convert content to a '[]string'
		}
	} else {
		mc.QueryString = append(mc.QueryString, *query)
	}

	mc.DrvConfig, err = drv.NewDefaultConfig(mc.InputConfig.Bucket, mc.InputConfig.Region, secret.AccessID, secret.SecretAccessKey)
	if err != nil {
		return Result{}, err
	}
	if isFlagPassed("b") {
		mc.InputConfig.Bucket = *bucket
		mc.DrvConfig.SetOutputBucket(mc.InputConfig.Bucket)
	}
	if isFlagPassed("d") {
		mc.InputConfig.Database = *database
	}
	if isFlagPassed("r") {
		mc.OutputConfig.Rowonly = *rowOnly
	}
	if isFlagPassed("m") {
		mc.OutputConfig.Moneywise = *moneyWise
	}
	if isFlagPassed("f") {
		mc.OutputConfig.Fastfail = *fastFail
	} else {
		mc.OutputConfig.Fastfail = true
	}
	if isFlagPassed("a") {
		mc.InputConfig.Admin = *admin
	}
	if isFlagPassed("y") {
		mc.OutputConfig.Style = *style
	}
	if isFlagPassed("o") {
		mc.OutputConfig.Render = *format
	}
	if mc.OutputConfig.Moneywise {
		mc.DrvConfig.SetMoneyWise(true)
	}
	mc.DrvConfig.SetDB(mc.InputConfig.Database)
	if !mc.InputConfig.Admin {
		mc.DrvConfig.SetReadOnly(true)
	}
	if err != nil {
		return Result{}, err
	}
	return Result{
		MyConfig: mc,
	}, nil
}

func expand(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}

	usr, err := user.Current()
	if err != nil {
		return "/tmp/"
	}
	return filepath.Join(usr.HomeDir, path[1:])
}
