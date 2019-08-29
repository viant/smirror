package pubsub

import (
	"context"
	"encoding/base64"
	"google.golang.org/api/pubsub/v1"
	"smirror/msgbus"
)

type service struct {
	*pubsub.Service
}


func (s *service) Publish(ctx context.Context, topic string, data []byte) ([]string, error) {
	publish := pubsub.NewProjectsService(s.Service).Topics.Publish(topic, &pubsub.PublishRequest{
		Messages: []*pubsub.PubsubMessage{
			{
				Data: base64.StdEncoding.EncodeToString(data),
			},
		},
	})
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
		Service:   srv,
	}, nil
}
