package mirror

import (
	"sync"
	"time"
)

const (
	//StatusOK status ok
	StatusOK = "ok"
	//StatusError status error
	StatusError = "error"
	//StatusNoMatch status no match
	StatusNoMatch = "noMatch"
)

//Request represents a mirror request
type Request struct {
	URL string
}

//Response represents a response
type Response struct {
	DestURLs    []string
	TimeTakenMs int
	Status      string
	Error       string

	startTime time.Time
	mutext    *sync.Mutex
}

//AddURL adds url to dest urls
func (r *Response) AddURL(URL string) {
	r.mutext.Lock()
	defer r.mutext.Unlock()
	r.DestURLs = append(r.DestURLs, URL)
}

//NewRequest create a request
func NewRequest(URL string) *Request {
	return &Request{URL: URL}
}

//NewResponse returns a new response
func NewResponse() *Response {
	return &Response{
		Status:    StatusOK,
		startTime: time.Now(),
		DestURLs:  make([]string, 0),
		mutext:    &sync.Mutex{},
	}
}
