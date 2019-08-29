package msgbus

import "context"

type Service interface {

	Publish(ctx context.Context, topic string, data []byte) ([]string, error)

}


