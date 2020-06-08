package athenadriver

// QIDMetaData is the meta data for QID
type QIDMetaData struct {
	QID         string
	dataScanned int64
	timestamp   int64
}

// AthenaCache is for Cached Query
type AthenaCache interface {
	// SetQID is to put query -> QIDMetaData into cache
	SetQID(query string, data QIDMetaData)

	// GetQID is to get QIDMetaData from cache by query string
	GetQID(query string) QIDMetaData

	// GetQuery is to get query string from cache by QID
	GetQuery(QID string) string
}
