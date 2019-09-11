package msgbus

import "context"

//Service represents message bus service
type Service interface {
	//Publish publishes data to message bus
	Publish(ctx context.Context, topic string, data []byte, attributes map[string]interface{}) ([]string, error)
}
