package tools


import (
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"errors"
)


/// DB MAP
/// a collection for a name
/// collection construct from requestinfo 


var ErrPostOne = errors.New("post one error")
var ErrGetStat = errors.New("get stat error")


// HTTPLogger holds api server's url
type RequestLogger struct {
	PostUrl string
	GetUrl string
}


// ErrInfo
type ErrInfo struct {
	Err int `json:"err"`
	Msg string `json:"msg"`
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


// RequestToInfo makes info from request
func RequestToInfo(req *http.Request, t time.Time) RequestInfo {
	var ri RequestInfo
	ri.Method = req.Method
	ri.Path = req.URL.Path
	ri.ClientIP = req.RemoteAddr
	ri.Time = t
	return ri
}


// NewRequestLogger gives a new instance
func NewRequestLogger(post, get string) *RequestLogger {
	return &RequestLogger{PostUrl: post, GetUrl: get}
}


// PostOne post one request log to  server
func (hl *RequestLogger) PostOne(req *http.Request) error {
	ri := RequestToInfo(req, time.Now())
	bt, _ := json.Marshal(ri)
	buf := bytes.NewBuffer(bt)
	resp, err := http.Post(hl.PostUrl, "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bt, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var ret ErrInfo
	err = json.Unmarshal(bt, &ret)
	if err != nil {
		return err
	}
	if ret.Err != 0 {
		return ErrPostOne
	}
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
