package contract

import "time"

//Request represents a mirror request
type Request struct {
	URL       string
	Attempt   int
	Timestamp time.Time
}

//NewRequest create a request
func NewRequest(URL string) *Request {
	return &Request{URL: URL, Timestamp: time.Now()}
}
