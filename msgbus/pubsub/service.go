package pubsub

import (
	"context"
	"encoding/base64"
	"fmt"
	"google.golang.org/api/pubsub/v1"
	"smirror/msgbus"
)

type service struct {
	*pubsub.Service
}

func (s *service) Publish(ctx context.Context, topic string, data []byte, attributes map[string]interface{}) ([]string, error) {
	request := &pubsub.PublishRequest{
		Messages: []*pubsub.PubsubMessage{
			{
				Data: base64.StdEncoding.EncodeToString(data),
			},
		},
	}
	if len(attributes) > 0 {
		request.Messages[0].Attributes = make(map[string]string)
		for k, v := range attributes {
			request.Messages[0].Attributes[k] = fmt.Sprintf("%s", v)
		}
	}
	publish := pubsub.NewProjectsService(s.Service).Topics.Publish(topic, request)

	publish.Context(ctx)
	response, err := publish.Do()
	if err != nil {
		return nil, err
	}
	return response.MessageIds, nil
}

func New(ctx context.Context) (msgbus.Service, error) {
	srv, err := pubsub.NewService(ctx)
	if err != nil {
		return nil, err
	}
	return &service{
		Service: srv,
	}, nil
}
