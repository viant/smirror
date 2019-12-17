package msg

import "smirror/base"

//Response represents a response
type Response struct {
	Status  string
	Error   string
	EventID string
	Size    int
	URL     string
}

//NewResponse creates a response
func NewResponse(eventID string) *Response {
	return &Response{
		EventID: eventID,
		Status:  base.StatusOK,
	}
}
