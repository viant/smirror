package msgbus

import "context"

//Service represents message bus service
type Service interface {
	//Publish publishes data to message bus
	Publish(ctx context.Context, request *Request) (*Response, error)
}

//Request represents request
type Request struct {
	Dest       string
	Data       []byte
	Attributes map[string]interface{}
}

//Response represents response
type Response struct {
	MessageIDs []string
}
