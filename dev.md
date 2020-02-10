- To generate UML graph

```go
goplantuml /opt/share/go/path/src/github.com/uber/athenasql/go \
> athenasql.puml
```

- Lint

```go
golangci-lint run
```

- fmt

```go
find go/ -type f -iregex '.*\.go' -exec go fmt '{}' +
```

- test

```go
✔ /opt/share/go/path/src/github.com/uber/athenasql/go [master|✚ 4…9] 
23:23 $ go test -v
=== RUN   TestAthenaConfig
--- PASS: TestAthenaConfig (0.00s)
=== RUN   TestAthenaConfigWrongS3Bucket
--- PASS: TestAthenaConfigWrongS3Bucket (0.00s)
=== RUN   TestAthenaConfigWrongRegion
--- PASS: TestAthenaConfigWrongRegion (0.00s)
=== RUN   TestAthenaConfigWrongWG
--- PASS: TestAthenaConfigWrongWG (0.00s)
=== RUN   TestAthenaConfigSafeString
--- PASS: TestAthenaConfigSafeString (0.00s)
=== RUN   TestReadOnlyCTAS1
--- PASS: TestReadOnlyCTAS1 (0.00s)
=== RUN   TestReadOnlyCTAS2
--- PASS: TestReadOnlyCTAS2 (0.00s)
=== RUN   TestReadOnlyCTAS3
--- PASS: TestReadOnlyCTAS3 (0.00s)
=== RUN   TestReadOnlyDROP
--- PASS: TestReadOnlyDROP (0.00s)
=== RUN   TestConnection_Prepare
--- PASS: TestConnection_Prepare (0.00s)
=== RUN   TestConnection_Begin
--- PASS: TestConnection_Begin (0.00s)
=== RUN   TestConnection_Close
--- PASS: TestConnection_Close (0.00s)
=== RUN   TestConnection_QueryContext
--- PASS: TestConnection_QueryContext (0.00s)
=== RUN   TestConnection_BeginTx
--- PASS: TestConnection_BeginTx (0.00s)
=== RUN   TestConnection_Transaction
--- PASS: TestConnection_Transaction (0.00s)
=== RUN   TestConnection_InterpolateParams
--- PASS: TestConnection_InterpolateParams (0.00s)
=== RUN   TestInterpolateParamsTooManyPlaceholders
--- PASS: TestInterpolateParamsTooManyPlaceholders (0.00s)
=== RUN   TestInterpolateParamsPlaceholderInString
--- PASS: TestInterpolateParamsPlaceholderInString (0.00s)
=== RUN   TestInterpolateParamsUint64
--- PASS: TestInterpolateParamsUint64 (0.00s)
=== RUN   TestCheckNamedValue
--- PASS: TestCheckNamedValue (0.00s)
=== RUN   TestConnector
--- PASS: TestConnector (0.00s)
=== RUN   TestDriver
--- PASS: TestDriver (0.00s)
=== RUN   TestObservability_Config
--- PASS: TestObservability_Config (0.00s)
=== RUN   TestObservability_Scope
--- PASS: TestObservability_Scope (0.00s)
=== RUN   TestObservability_Logger
--- PASS: TestObservability_Logger (0.00s)
=== RUN   TestObservability_Log
--- PASS: TestObservability_Log (0.00s)
=== RUN   TestOnePageSuccess
--- PASS: TestOnePageSuccess (0.00s)
=== RUN   TestNextFailure
--- PASS: TestNextFailure (0.00s)
=== RUN   TestMultiplePages
--- PASS: TestMultiplePages (0.00s)
PASS
ok  	github.com/uber/athenasql/go	0.005s
```



```go
$go test -coverprofile=coverage.out github.com/uber/athenasql/go  && go tool cover -func=coverage.out|grep -v 100.0%
ok  	github.com/uber/athenasql/go	9.241s	coverage: 100.0% of statements
```
