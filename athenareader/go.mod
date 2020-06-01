module github.com/uber/athenadriver/athenareader

go 1.13

replace github.com/uber/athenadriver => /opt/share/go/path/src/github.com/uber/athenadriver

require (
	github.com/uber/athenadriver v1.1.6
	go.uber.org/config v1.4.0
	go.uber.org/fx v1.12.0
)
