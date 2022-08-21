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

// Package athenadriver is a fully-featured Go database/sql driver for
// Amazon AWS Athena developed at Uber Technologies Inc.
//
// It provides a hassle-free way of querying AWS Athena database with Go
// standard library. It not only provides basic features of Athena Go SDK, but
// addresses some of its limitation, improves and extends it.Except the basic
// features provided by Go database/sql like error handling, database pool
// and reconnection, athenadriver supports the following features out of box:
//
//   - Support multiple AWS authorization methods
//   - Full support of Athena Basic Data Types
//   - Full support of Athena Advanced Type for queries with Geospatial identifiers, ML and UDFs
//   - Full support of ALL Athena Query Statements, including DDL, DML and UTILITY
//   - Support newly added INSERT INTO...VALUES
//   - Full support of Athena Basic Data Types
//   - Athena workgroup and tagging support including remote workgroup creation
//   - Go sql's Prepared statement support
//   - Go sql's DB.Exec() and db.ExecContext() support
//   - Query cancelling support
//   - Mask columns with specific values
//   - Database missing value handling
//   - Read-Only mode
//
// Amazon Athena is an interactive query service that lets you use standard
// SQL to analyze data directly in Amazon S3. You can point Athena at your data
// in Amazon S3 and run ad-hoc queries and get results in seconds. Athena is
// serverless, so there is no infrastructure to set up or manage. You pay only
// for the queries you run. Athena scales automatically—executing queries
// in parallel—so results are fast, even with large datasets and complex queries.
// Author: Henry Fuheng Wu (wufuheng@gmail.com, henry.wu@uber.com)
package athenadriver
