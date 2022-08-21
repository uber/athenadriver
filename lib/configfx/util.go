package configfx

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func setUpFlagUsage(context.Context) error {
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
			"AUTHOR\n\tHenry Fuheng Wu (wufuheng@gmail.com)\n\n" +
			"REPORTING BUGS\n\thttps://github.com/uber/athenadriver\n"
		fmt.Fprintf(commandLine.Output(), preBody)
		fmt.Fprintf(commandLine.Output(),
			"SYNOPSIS\n\n\t%s [-v] [-b OUTPUT_BUCKET] [-d DATABASE_NAME] [-q QUERY_STRING_OR_FILE] [-r] [-a] [-m] [-y STYLE_NAME] [-o OUTPUT_FORMAT]\n\nDESCRIPTION\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(commandLine.Output(), desc)
	}
	return nil
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

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
