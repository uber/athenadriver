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
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/aws/aws-sdk-go/service/athena/athenaiface"
)

// genQueryResultsOutputByToken is a function type with string as parameter.
type genQueryResultsOutputByToken func(token string) (*athena.GetQueryResultsOutput, error)

// mockAthenaClient is a type embedding of athenaiface.AthenaAPI, which we want to mock against.
type mockAthenaClient struct {
	// inner type is interface athenaiface.AthenaAPI
	athenaiface.AthenaAPI

	// queryToResultsGenMap is a map from string to a function type genQueryResultsOutputByToken.
	queryToResultsGenMap map[string]genQueryResultsOutputByToken

	CreateWGStatus bool
	GetWGStatus    bool
	WGDisabled     bool
}

func newMockAthenaClient() *mockAthenaClient {
	var m = mockAthenaClient{
		queryToResultsGenMap: map[string]genQueryResultsOutputByToken{
			"SELECT_OK":                            MultiplePagesQueryResponse,
			"SELECT_GetQueryResults_ERR":           MultiplePagesQueryFailedResponse,
			"SELECT_EMPTY_ROW_IN_PAGE":             MultiplePagesEmptyRowInPageResponse,
			"show":                                 ShowResponse,
			"RowsNextFailed":                       NextFailedResponse,
			"1coloumn0row":                         OneColumnZeroRowResponse,
			"1coloumn0row_valid":                   OneColumnZeroRowResponseValid,
			"column_more_than_row_fields":          ColumnMoreThanRowFieldResponse,
			"row_fields_more_than_column":          RowFieldMoreThanColumnsResponse,
			"missing_data_resp":                    MissingDataResponse,
			"missing_data_resp2":                   headPageWithColumnButNoRowResponse,
			"PING_OK_QID":                          PingResponse,
			"SELECTExecContext_OK_QID":             PingResponse,
			"SELECTQueryContext_OK_QID":            PingResponse,
			"00000000-0000-0000-0000-000000000000": PingResponse,
			"pc:get_query_id":                      PingResponse,
			"FAILED_AFTER_GETQID":                  MissingDataResponse,
		},
	}
	return &m
}

// GetQueryResults is a mock against athenaiface.AthenaAPI.GetQueryResults().
func (m *mockAthenaClient) GetQueryResults(query *athena.
	GetQueryResultsInput) (*athena.GetQueryResultsOutput, error) {
	var nextToken = ""
	if query.NextToken != nil {
		nextToken = *query.NextToken
	}
	return m.queryToResultsGenMap[*query.QueryExecutionId](nextToken)
}

// GetQueryResultsWithContext is a mock against athenaiface.AthenaAPI.GetQueryResultsWithContext().
func (m *mockAthenaClient) GetQueryResultsWithContext(ctx aws.Context,
	query *athena.GetQueryResultsInput, opt ...request.Option) (*athena.
	GetQueryResultsOutput, error) {
	var nextToken = ""
	if query.NextToken != nil {
		nextToken = *query.NextToken
	}
	if *query.QueryExecutionId == "GetQueryResultsWithContext_return_error" {
		return nil, ErrTestMockGeneric
	}
	if nextToken == "GetQueryResultsWithContext_return_error" {
		return nil, ErrTestMockGeneric
	}
	return m.queryToResultsGenMap[*query.QueryExecutionId](nextToken)
}

func (m *mockAthenaClient) GetWorkGroupWithContext(ctx aws.Context, gwi *athena.GetWorkGroupInput,
	opt ...request.Option) (*athena.GetWorkGroupOutput, error) {
	if m.GetWGStatus {
		enabled := "ENABLED"
		if m.WGDisabled {
			enabled = "DISABLED"
		}
		w := athena.WorkGroup{
			State: &enabled,
		}
		a := athena.GetWorkGroupOutput{
			WorkGroup: &w,
		}
		return &a, nil
	}
	return nil, ErrTestMockGeneric
}

func (m *mockAthenaClient) CreateWorkGroup(*athena.CreateWorkGroupInput) (
	*athena.CreateWorkGroupOutput, error) {
	if !m.CreateWGStatus {
		return nil, ErrTestMockGeneric
	}
	a := athena.CreateWorkGroupOutput{}
	return &a, nil
}

func (m *mockAthenaClient) StartQueryExecution(s *athena.
	StartQueryExecutionInput) (*athena.StartQueryExecutionOutput, error) {
	if strings.ToLower(*s.QueryString) == "select 1" { // Ping
		qid := "PING_OK_QID"
		return &athena.StartQueryExecutionOutput{
			QueryExecutionId: &qid,
		}, nil
	}
	if *s.QueryString == "SELECTExecContext_OK" { // Ping
		qid := "SELECTExecContext_OK_QID"
		return &athena.StartQueryExecutionOutput{
			QueryExecutionId: &qid,
		}, nil
	}
	if *s.QueryString == "SELECTQueryContext_OK" ||
		*s.QueryString == "SELECTQueryContext_'OK'" { // Ping
		qid := "SELECTQueryContext_OK_QID"
		return &athena.StartQueryExecutionOutput{
			QueryExecutionId: &qid,
		}, nil
	}
	if *s.QueryString == "SELECTQueryContext_CANCEL_OK" { // Ping
		qid := "SELECTQueryContext_CANCEL_OK_QID"
		return &athena.StartQueryExecutionOutput{
			QueryExecutionId: &qid,
		}, nil
	}
	if *s.QueryString == "SELECTQueryContext_AWS_CANCEL" { // Ping
		qid := "SELECTQueryContext_AWS_CANCEL_QID"
		return &athena.StartQueryExecutionOutput{
			QueryExecutionId: &qid,
		}, nil
	}
	if *s.QueryString == "SELECTQueryContext_AWS_FAIL" { // Ping
		qid := "SELECTQueryContext_AWS_FAIL_QID"
		return &athena.StartQueryExecutionOutput{
			QueryExecutionId: &qid,
		}, nil
	}
	if *s.QueryString == "SELECTQueryContext_CANCEL_FAIL" { // Ping
		qid := "SELECTQueryContext_CANCEL_FAIL_QID"
		return &athena.StartQueryExecutionOutput{
			QueryExecutionId: &qid,
		}, nil
	}
	if *s.QueryString == "SELECTQueryContext_TIMEOUT" { // Ping
		qid := "SELECTQueryContext_TIMEOUT_QID"
		return &athena.StartQueryExecutionOutput{
			QueryExecutionId: &qid,
		}, nil
	}
	if *s.QueryString == "StartQueryExecution_nil_error" {
		return nil, ErrTestMockGeneric
	}
	if *s.QueryString == "When_StartQueryExecution_Succeed_but_GetQueryExecutionWithContext_return_nil_and_error" {
		qid := "When_StartQueryExecution_Succeed_but_GetQueryExecutionWithContext_return_nil_and_error_QID"
		return &athena.StartQueryExecutionOutput{
			QueryExecutionId: &qid,
		}, nil
	}
	if *s.QueryString == "StartQueryExecution_OK_GetQueryExecutionWithContext_QueryExecutionStateCancelled" {
		qid := "QueryExecutionStateCancelled_QID"
		return &athena.StartQueryExecutionOutput{
			QueryExecutionId: &qid,
		}, nil
	}
	if *s.QueryString == "StartQueryExecution_OK_GetQueryExecutionWithContext_QueryExecutionStateFailed" {
		qid := "QueryExecutionStateFailed_QID"
		return &athena.StartQueryExecutionOutput{
			QueryExecutionId: &qid,
		}, nil
	}
	if *s.QueryString == "FAILED_AFTER_GETQID" {
		qid := "FAILED_AFTER_GETQID_123"
		return &athena.StartQueryExecutionOutput{
			QueryExecutionId: &qid,
		}, fmt.Errorf("FAILED_AFTER_GETQID_FAILED")
	}
	if *s.QueryString == "FAILED_AFTER_GETQID2" {
		qid := "FAILED_AFTER_GETQID_123"
		return &athena.StartQueryExecutionOutput{
				QueryExecutionId: &qid,
			}, awserr.NewRequestFailure(awserr.New("a", "b", fmt.Errorf("FAILED_AFTER_GETQID_FAILED")),
				100, qid)
	}
	return nil, nil
}

func (m *mockAthenaClient) GetQueryExecutionWithContext(c aws.Context,
	input *athena.GetQueryExecutionInput, o ...request.Option) (*athena.GetQueryExecutionOutput, error) {
	if *input.QueryExecutionId == "When_StartQueryExecution_Succeed_but_GetQueryExecutionWithContext_return_nil_and_error_QID" {
		return nil, ErrTestMockGeneric
	}
	if *input.QueryExecutionId == "QueryExecutionStateCancelled_QID" {
		return nil, context.Canceled
	}
	if *input.QueryExecutionId == "QueryExecutionStateFailed_QID" {
		return nil, ErrTestMockFailedByAthena
	}
	if *input.QueryExecutionId == "PING_OK_QID" {
		ping := "PING_OK_QID"
		stat := athena.QueryExecutionStateSucceeded
		return &athena.GetQueryExecutionOutput{
			QueryExecution: &athena.QueryExecution{
				Query:            &ping,
				QueryExecutionId: &ping,
				Status: &athena.QueryExecutionStatus{
					State: &stat,
				},
			},
		}, nil
	}
	if *input.QueryExecutionId == "SELECTExecContext_OK_QID" {
		ping := "SELECTExecContext_OK_QID"
		stat := athena.QueryExecutionStateSucceeded
		var dataScanned = int64(123)
		return &athena.GetQueryExecutionOutput{
			QueryExecution: &athena.QueryExecution{
				Query:            &ping,
				QueryExecutionId: &ping,
				Status: &athena.QueryExecutionStatus{
					State: &stat,
				},
				Statistics: &athena.QueryExecutionStatistics{
					DataScannedInBytes: &dataScanned,
				},
			},
		}, nil
	}
	if *input.QueryExecutionId == "SELECTQueryContext_OK_QID" {
		ping := "SELECTQueryContext_OK_QID"
		stat := athena.QueryExecutionStateSucceeded
		stt := "DDL"
		return &athena.GetQueryExecutionOutput{
			QueryExecution: &athena.QueryExecution{
				Query:            &ping,
				QueryExecutionId: &ping,
				Status: &athena.QueryExecutionStatus{
					State: &stat,
				},
				StatementType: &stt,
			},
		}, nil
	}
	if *input.QueryExecutionId == "SELECTQueryContext_CANCEL_OK_QID" {
		ping := "SELECTQueryContext_CANCEL_OK_QID"
		stat := athena.QueryExecutionStateQueued
		stt := "DDL"
		var dataScanned = int64(123)
		return &athena.GetQueryExecutionOutput{
			QueryExecution: &athena.QueryExecution{
				Query:            &ping,
				QueryExecutionId: &ping,
				Status: &athena.QueryExecutionStatus{
					State: &stat,
				},
				StatementType: &stt,
				Statistics: &athena.QueryExecutionStatistics{
					DataScannedInBytes: &dataScanned,
				},
			},
		}, nil
	}
	if *input.QueryExecutionId == "SELECTQueryContext_AWS_CANCEL_QID" {
		ping := "SELECTQueryContext_AWS_CANCEL_QID"
		var dataScanned = int64(123)
		stat := athena.QueryExecutionStateCancelled
		return &athena.GetQueryExecutionOutput{
			QueryExecution: &athena.QueryExecution{
				Query:            &ping,
				QueryExecutionId: &ping,
				Status: &athena.QueryExecutionStatus{
					State: &stat,
				},
				Statistics: &athena.QueryExecutionStatistics{
					DataScannedInBytes: &dataScanned,
				},
			},
		}, nil
	}
	if *input.QueryExecutionId == "SELECTQueryContext_AWS_FAIL_QID" {
		ping := "SELECTQueryContext_AWS_FAIL_QID"
		stat := athena.QueryExecutionStateFailed
		reason := "something_broken"
		return &athena.GetQueryExecutionOutput{
			QueryExecution: &athena.QueryExecution{
				Query:            &ping,
				QueryExecutionId: &ping,
				Status: &athena.QueryExecutionStatus{
					State:             &stat,
					StateChangeReason: &reason,
				},
			},
		}, nil
	}
	if *input.QueryExecutionId == "SELECTQueryContext_CANCEL_FAIL_QID" {
		ping := "SELECTQueryContext_CANCEL_FAIL_QID"
		stat := athena.QueryExecutionStateQueued
		return &athena.GetQueryExecutionOutput{
			QueryExecution: &athena.QueryExecution{
				Query:            &ping,
				QueryExecutionId: &ping,
				Status: &athena.QueryExecutionStatus{
					State: &stat,
				},
			},
		}, nil
	}
	if *input.QueryExecutionId == "SELECTQueryContext_TIMEOUT_QID" {
		ping := "SELECTQueryContext_TIMEOUT_QID"
		stat := athena.QueryExecutionStateQueued
		stt := "TIMEOUT_NOW"
		return &athena.GetQueryExecutionOutput{
			QueryExecution: &athena.QueryExecution{
				Query:            &ping,
				QueryExecutionId: &ping,
				Status: &athena.QueryExecutionStatus{
					State: &stat,
				},
				StatementType: &stt,
			},
		}, nil
	}
	if *input.QueryExecutionId == "c89088ab-595d-4ee6-a9ce-73b55aeb8900" {
		ping := "SELECTQueryContext_CANCEL_OK_QID"
		stat := athena.QueryExecutionStateQueued
		stt := "DDL"
		var dataScanned = int64(123)
		return &athena.GetQueryExecutionOutput{
			QueryExecution: &athena.QueryExecution{
				Query:            &ping,
				QueryExecutionId: &ping,
				Status: &athena.QueryExecutionStatus{
					State: &stat,
				},
				StatementType: &stt,
				Statistics: &athena.QueryExecutionStatistics{
					DataScannedInBytes: &dataScanned,
				},
			},
		}, nil
	}
	return nil, ErrTestMockGeneric
}

func (m *mockAthenaClient) StopQueryExecutionWithContext(ctx aws.Context, input *athena.StopQueryExecutionInput,
	opt ...request.Option) (*athena.StopQueryExecutionOutput, error) {
	if *input.QueryExecutionId == "SELECTQueryContext_CANCEL_OK_QID" {
		return &athena.StopQueryExecutionOutput{}, nil
	}
	if *input.QueryExecutionId == "SELECTQueryContext_CANCEL_FAIL_QID" {
		return nil, ErrTestMockGeneric
	}
	if *input.QueryExecutionId == "c89088ab-595d-4ee6-a9ce-73b55aeb8954" {
		return &athena.StopQueryExecutionOutput{}, nil
	}
	if *input.QueryExecutionId == "c89088ab-595d-4ee6-a9ce-73b55aeb8955" {
		return nil, ErrTestMockGeneric
	}
	return nil, ErrTestMockGeneric
}

func MultiplePagesQueryResponse(token string) (*athena.GetQueryResultsOutput, error) {
	columns := createTestColumns()
	switch token {
	case "":
		var nextToken = "a1"
		return newRandomHeaderResultPage(columns, &nextToken, 6), nil
	case "a1":
		nextToken := "a2"
		return newRandomHeaderlessResultPage(columns, &nextToken, 10), nil
	case "a2":
		nextToken := "a3"
		return newRandomHeaderlessResultPage(columns, &nextToken, 5), nil
	case "a3":
		nextToken := "a4"
		return newRandomHeaderlessResultPage(columns, &nextToken, 5), nil
	case "a4":
		return newRandomHeaderlessResultPage(columns, nil, 10), nil
	default:
		return nil, ErrTestMockGeneric
	}
}

// page contains 0 row, missing row in the page
func MultiplePagesQueryFailedResponse(token string) (*athena.GetQueryResultsOutput,
	error) {
	columns := createTestColumns()
	switch token {
	case "":
		var nextToken = "a1"
		return newRandomHeaderResultPage(columns, &nextToken, 6), nil
	case "a1":
		nextToken := "a2"
		return newRandomHeaderlessResultPage(columns, &nextToken, 3), nil
	case "a2":
		nextToken := "a3"
		return newRandomHeaderlessResultPage(columns, &nextToken, 5), nil
	case "a3":
		nextToken := "GetQueryResultsWithContext_return_error"
		return newRandomHeaderlessResultPage(columns, &nextToken, 5), nil
	case "a4":
		return newRandomHeaderlessResultPage(columns, nil, 10), nil
	default:
		return nil, ErrTestMockGeneric
	}
}

// page contains 0 row, missing row in the page
func MultiplePagesEmptyRowInPageResponse(token string) (*athena.GetQueryResultsOutput,
	error) {
	columns := createTestColumns()
	switch token {
	case "":
		var nextToken = "a1"
		return newRandomHeaderResultPage(columns, &nextToken, 6), nil
	case "a1":
		nextToken := "a2"
		return newRandomHeaderlessResultPage(columns, &nextToken, 0), nil
	case "a2":
		nextToken := "a3"
		return newRandomHeaderlessResultPage(columns, &nextToken, 5), nil
	case "a3":
		nextToken := "GetQueryResultsWithContext_return_error"
		return newRandomHeaderlessResultPage(columns, &nextToken, 5), nil
	case "a4":
		return newRandomHeaderlessResultPage(columns, nil, 10), nil
	default:
		return nil, ErrTestMockGeneric
	}
}

func ShowResponse(_ string) (*athena.GetQueryResultsOutput, error) {
	columns := []*athena.ColumnInfo{
		newColumnInfo("partition", "string"),
	}
	return newRandomHeaderResultPage(columns, nil, 6), nil
}

func OneColumnZeroRowResponse(token string) (*athena.GetQueryResultsOutput,
	error) {
	switch token {
	case "":
		c := newColumnInfo("a", nil)
		getQueryResultsOutput := &athena.GetQueryResultsOutput{
			ResultSet: &athena.ResultSet{
				ResultSetMetadata: &athena.ResultSetMetadata{
					ColumnInfo: []*athena.ColumnInfo{
						c,
					},
				},
			},
		}
		return getQueryResultsOutput, nil
	default:
		return nil, ErrTestMockGeneric
	}
}

func OneColumnZeroRowResponseValid(token string) (*athena.GetQueryResultsOutput,
	error) {
	switch token {
	case "":
		c := newColumnInfo("rows", nil)
		var i int64 = 1024
		getQueryResultsOutput := &athena.GetQueryResultsOutput{
			ResultSet: &athena.ResultSet{
				ResultSetMetadata: &athena.ResultSetMetadata{
					ColumnInfo: []*athena.ColumnInfo{
						c,
					},
				},
			},
			UpdateCount: &i,
		}
		return getQueryResultsOutput, nil
	default:
		return nil, ErrTestMockGeneric
	}
}

func ColumnMoreThanRowFieldResponse(token string) (*athena.GetQueryResultsOutput,
	error) {
	switch token {
	case "":
		c1 := newColumnInfo("c1", nil)
		c2 := newColumnInfo("c2", nil)
		var i int64 = 1024
		getQueryResultsOutput := &athena.GetQueryResultsOutput{
			ResultSet: &athena.ResultSet{
				ResultSetMetadata: &athena.ResultSetMetadata{
					ColumnInfo: []*athena.ColumnInfo{
						c1, c2,
					},
				},
				Rows: []*athena.Row{
					randRow([]*athena.ColumnInfo{
						c1,
					}),
				},
			},
			UpdateCount: &i,
		}
		return getQueryResultsOutput, nil
	default:
		return nil, ErrTestMockGeneric
	}
}

func RowFieldMoreThanColumnsResponse(token string) (*athena.GetQueryResultsOutput,
	error) {
	switch token {
	case "":
		c1 := newColumnInfo("c1", nil)
		c2 := newColumnInfo("c2", nil)
		var i int64 = 1024
		getQueryResultsOutput := &athena.GetQueryResultsOutput{
			ResultSet: &athena.ResultSet{
				ResultSetMetadata: &athena.ResultSetMetadata{
					ColumnInfo: []*athena.ColumnInfo{
						c1,
					},
				},
				Rows: []*athena.Row{
					randRow([]*athena.ColumnInfo{
						c1, c2,
					}),
				},
			},
			UpdateCount: &i,
		}
		return getQueryResultsOutput, nil
	default:
		return nil, ErrTestMockGeneric
	}
}

func MissingDataResponse(token string) (*athena.GetQueryResultsOutput,
	error) {
	switch token {
	case "":
		c1 := newColumnInfo("c1", "integer")
		var i int64 = 1024
		getQueryResultsOutput := &athena.GetQueryResultsOutput{
			ResultSet: &athena.ResultSet{
				ResultSetMetadata: &athena.ResultSetMetadata{
					ColumnInfo: []*athena.ColumnInfo{
						c1,
					},
				},
				Rows: []*athena.Row{
					missingDataRow([]*athena.ColumnInfo{
						c1,
					}),
				},
			},
			UpdateCount: &i,
		}
		return getQueryResultsOutput, nil
	default:
		return nil, ErrTestMockGeneric
	}
}

func headPageWithColumnButNoRowResponse(token string) (*athena.GetQueryResultsOutput,
	error) {
	switch token {
	case "":
		c2 := newColumnInfo("c2", "string")
		var i int64 = 1024
		getQueryResultsOutput := &athena.GetQueryResultsOutput{
			ResultSet: &athena.ResultSet{
				ResultSetMetadata: &athena.ResultSetMetadata{
					ColumnInfo: []*athena.ColumnInfo{
						c2,
					},
				},
				Rows: []*athena.Row{
					missingDataRow([]*athena.ColumnInfo{
						c2,
					}),
				},
			},
			UpdateCount: &i,
		}
		return getQueryResultsOutput, nil
	default:
		return nil, ErrTestMockGeneric
	}
}

func PingResponse(token string) (*athena.GetQueryResultsOutput,
	error) {
	switch token {
	case "":
		c2 := newColumnInfo("_col0", "integer")
		var i int64 = 1024
		getQueryResultsOutput := &athena.GetQueryResultsOutput{
			ResultSet: &athena.ResultSet{
				ResultSetMetadata: &athena.ResultSetMetadata{
					ColumnInfo: []*athena.ColumnInfo{
						c2,
					},
				},
				Rows: []*athena.Row{
					randRow([]*athena.ColumnInfo{
						c2,
					}),
				},
			},
			UpdateCount: &i,
		}
		return getQueryResultsOutput, nil
	default:
		return nil, ErrTestMockGeneric
	}
}

func NextFailedResponse(token string) (*athena.GetQueryResultsOutput, error) {
	columns := createTestColumns()
	switch token {
	case "":
		var nextToken = "p1"
		return newRandomHeaderResultPage(columns, &nextToken, 5), nil
	default:
		return nil, ErrTestMockGeneric
	}
}

func createTestColumns() []*athena.ColumnInfo {
	return []*athena.ColumnInfo{
		newColumnInfo("test_array", "array"),
		newColumnInfo("active", "boolean"),
		newColumnInfo("company_name", "string"),
		newColumnInfo("project", "string"),
		newColumnInfo("uid", "integer"),
		newColumnInfo("regitser_date", "date"),
		newColumnInfo("regitser_ts", "timestamp"),
	}
}
