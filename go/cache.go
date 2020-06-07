package athenadriver

type QIDMetaData struct {
	QID         string
	dataScanned int64
	timestamp   int64
}

type AthenaCache interface {
	// SetQID is to put query -> QIDMetaData into cache
	SetQID(query string, data QIDMetaData)
}
