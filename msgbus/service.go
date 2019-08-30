package msgbus

import "context"

type Service interface {
	Publish(ctx context.Context, topic string, data []byte, attributes map[string]interface{}) ([]string, error)
}
