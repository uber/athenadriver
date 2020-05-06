package main

import (
	"database/sql"
	"os"

	secret "github.com/uber/athenadriver/examples/constants"
	drv "github.com/uber/athenadriver/go"
)

func main() {
	// 1. Set AWS Credential in Driver Config.
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	conf, err := drv.NewDefaultConfig(secret.OutputBucket, secret.Region, secret.AccessID, secret.SecretAccessKey)
	conf.SetLogging(true)
	if err != nil {
		panic(err)
		return
	}

	// 2. Open Connection.
	dsn := conf.Stringify()
	db, _ := sql.Open(drv.DriverName, dsn)

	// 3. Query with pseudo command `pc:get_query_id`
	var qidStatus string
	_ = db.QueryRow("pc:stop_query_id c89088ab-595d-4ee6-a9ce-73b55aeb8953").Scan(&qidStatus)
	println("Stop Query ID c89088ab-595d-4ee6-a9ce-73b55aeb8953 returns:", qidStatus)
}

/*
Sample Output:
Stop Query ID c89088ab-595d-4ee6-a9ce-73b55aeb8953 returns: OK
*/
