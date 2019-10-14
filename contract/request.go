package contract

//Request represents a mirror request
type Request struct {
	URL string
}

//NewRequest create a request
func NewRequest(URL string) *Request {
	return &Request{URL: URL}
}
