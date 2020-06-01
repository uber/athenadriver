// Copyright (c) 2020 Uber Technologies, Inc.
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
	"flag"
	"fmt"
	secret "github.com/uber/athenadriver/examples/constants"
	drv "github.com/uber/athenadriver/go"
	"go.uber.org/config"
	"go.uber.org/fx"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Module is to provide dependency of Configuration to main app
var Module = fx.Provide(new)

// Params defines the dependencies or inputs
type Params struct {
	fx.In
	Shutdowner fx.Shutdowner
}

// ReaderInputConfig is to represent the output section of configuration file
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
	QueryString string
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
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	var commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.Usage = func() {
		preBody := "NAME\n\tathenareader - read athena data from command line\n\n"
		desc := "\nEXAMPLES\n\n" +
			"\t$ athenareader -d sampledb -q \"select request_timestamp,elb_name from elb_logs limit 2\"\n" +
			"\trequest_timestamp,elb_name\n" +
			"\t2015-01-03T00:00:00.516940Z,elb_demo_004\n" +
			"\t2015-01-03T00:00:00.902953Z,elb_demo_004\n\n" +
			"\t$ athenareader -d sampledb -q \"select request_timestamp,elb_name from elb_logs limit 2\" -r\n" +
			"\t2015-01-05T20:00:01.206255Z,elb_demo_002\n" +
			"\t2015-01-05T20:00:01.612598Z,elb_demo_008\n\n" +
			"\t$ athenareader -d sampledb -b s3://my-athena-query-result -q tools/query.sql\n" +
			"\trequest_timestamp,elb_name\n" +
			"\t2015-01-06T00:00:00.516940Z,elb_demo_009\n\n" +
			"\n\tAdd '-m' to enable moneywise mode. The first line will display query cost under moneywise mode.\n\n" +
			"\t$ athenareader -b s3://athena-query-result -q 'select count(*) as cnt from sampledb.elb_logs' -m\n" +
			"\tquery cost: 0.00184898369752772851 USD\n" +
			"\tcnt\n" +
			"\t1356206\n\n" +
			"\n\tAdd '-a' to enable admin mode. Database write is enabled at driver level under admin mode.\n\n" +
			"\t$ athenareader -b s3://athena-query-result -q 'DROP TABLE IF EXISTS depreacted_table' -a\n" +
			"\t\n" +
			"AUTHOR\n\tHenry Fuheng Wu (henry.wu@uber.com)\n\n" +
			"REPORTING BUGS\n\thttps://github.com/uber/athenadriver\n"
		fmt.Fprintf(commandLine.Output(), preBody)
		fmt.Fprintf(commandLine.Output(),
			"SYNOPSIS\n\n\t%s [-v] [-b OUTPUT_BUCKET] [-d DATABASE_NAME] [-q QUERY_STRING_OR_FILE] [-r] [-a] [-m] [-y STYLE_NAME] [-o OUTPUT_FORMAT]\n\nDESCRIPTION\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(commandLine.Output(), desc)
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func new(p Params) (Result, error) {
	var bucket = flag.String("b", secret.OutputBucket, "Athena resultset output bucket")
	var database = flag.String("d", "default", "The database you want to query")
	var query = flag.String("q", "select 1", "The SQL query string or a file containing SQL string")
	var rowOnly = flag.Bool("r", false, "Display rows only, don't show the first row as columninfo")
	var moneyWise = flag.Bool("m", false, "Enable moneywise mode to display the query cost as the first line of the output")
	var versionFlag = flag.Bool("v", false, "Print the current version and exit")
	var admin = flag.Bool("a", false, "Enable admin mode, so database write(create/drop) is allowed at athenadriver level")
	var style = flag.String("y", "default", "Output rendering style")
	var format = flag.String("o", "csv", "Output format(options: table, markdown, csv, html)")

	flag.Parse()
	switch {
	case *versionFlag:
		println("Current build version: v" + drv.DriverVersion)
		p.Shutdowner.Shutdown()
		os.Exit(0)
		return Result{}, fmt.Errorf("no")
	}

	var mc = AthenaDriverConfig{}
	var (
		provider *config.YAML
		err      error
	)

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
				p.Shutdowner.Shutdown()
				os.Exit(2)
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
				p.Shutdowner.Shutdown()
				os.Exit(2)
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

	mc.QueryString = *query
	if _, err := os.Stat(*query); err == nil {
		b, err := ioutil.ReadFile(*query)
		if err == nil {
			mc.QueryString = string(b) // convert content to a 'string'
		}
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

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func downloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
