package tools


import (
	"net/http"
	"time"
)


/// DB MAP
/// a collection for a name
/// collection construct from requestinfo 


// HTTPLogger holds api server's url
type RequestLogger struct {
	PostUrl string
	GetUrl string
}


// RequestLoggerStat holds logs from start to end
type RequestLoggerStat struct {
	MethodCount map[string]int64 `json:"methodCount"`
	PathCount map[string]int64 `json:"pathCount"`
	ClientIPCount map[string]int64 `json:"clientIPCount"`
}


// RequestInfo retrieves all info in http.Request
type RequestInfo struct {
	Method string `json:"method"` // request method, get, post, put, head etc.
	Path string `json:"path"` // url.Path
	ClientIP string `json:"clientIP"` // client ip
	Time time.Time `json:"time"` // when does the request fired
}


// NewRequestLogger gives a new instance
func NewRequestLogger(post, get string) *RequestLogger {
	return &RequestLogger{PostUrl: post, GetUrl: get}
}


// PostOne post one request log to  server
func (hl *RequestLogger) PostOne(req *http.Request) error {
	return nil
}


// GetStat get log stat from start to end
func (hl *RequestLogger) GetStat(start, end time.Time) (*RequestLoggerStat, error) {
	return nil, nil
}


// ServeRequestLogger is server side object of request logger
type ServeRequestLogger struct {
	DBUrl string
	Granularity time.Time
	Name string
}


// NewServeRequestLogger create instance of server side logger
func NewServeRequestLogger(dbu, name string, gra time.Time) *ServeRequestLogger {
	return &ServeRequestLogger{DBUrl: dbu, Granularity: gra, Name: name}
}


// OneRequest handle one request
// record it into mongodb
func (s *ServeRequestLogger) OneRequest (r *RequestInfo) error {
	return nil
}


// GetStat get data from mongodb
// and give back ip info
func (s *ServeRequestLogger) GetStat (start, end time.Time) (*RequestLoggerStat, error) {
	return nil, nil
}
