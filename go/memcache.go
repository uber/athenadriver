package athenadriver

type CacheInMem struct {
	QIDToDataSize map[string]int64
	QIDToTimeStamp map[string]int64
	queryToQID map[string]string
}

var AthenaCacheInMem = newCache()

func newCache() *CacheInMem {
	return &CacheInMem{
		QIDToDataSize: map[string]int64{},
		QIDToTimeStamp: map[string]int64{},
		queryToQID: map[string]string{},
	}
}

func GetQID(query string) string{
	if val, ok := AthenaCacheInMem.queryToQID[query]; ok {
		return val
	}
	return ""
}

func SetQID(query string, dataSize int64, QID string, t int64){
	AthenaCacheInMem.QIDToDataSize[QID] = dataSize
	AthenaCacheInMem.queryToQID[query] = QID
	AthenaCacheInMem.QIDToTimeStamp[QID] = t
}

