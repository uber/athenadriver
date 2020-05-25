package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	drv "github.com/uber/athenadriver/go"
	"os"
)

type response struct {
	QueryResult string `json:"result"`
}

// Make sure to select a role which can query Athena!
// https://epsagon.com/blog/getting-started-with-aws-lambda-and-go/
// https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 1. Set AWS Credential in Driver Config.
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	conf, err := drv.NewDefaultConfig("s3://athena-query-result/lambda/",
		drv.DummyRegion, drv.DummyAccessID, drv.DummySecretAccessKey)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// 2. Open Connection.
	db, err := sql.Open(drv.DriverName, conf.Stringify())
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string(err.Error()), StatusCode: 500}, err
	}
	// 3. Query
	rows, err := db.QueryContext(ctx, "select 123")
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string(err.Error()), StatusCode: 500}, err
	}
	defer rows.Close()
	resp := &response{
		QueryResult: drv.ColsRowsToCSV(rows),
	}
	body, err := json.Marshal(resp)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}
func main() {
	lambda.Start(handleRequest)
}
