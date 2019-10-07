package pubsub

import (
	"context"
	"encoding/base64"
	"fmt"
	"google.golang.org/api/pubsub/v1"
	"smirror/msgbus"
	"strings"
)

type service struct {
	*pubsub.Service
	projectID string
}

func (s *service) topicInProject(topic string) string {
	if strings.Count(topic, "/") > 0 {
		return topic
	}
	return fmt.Sprintf("projects/%s/topics/%s", s.projectID, topic)
}

//Publish publishes data to message bus
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

	topic = s.topicInProject(topic)
	publish := pubsub.NewProjectsService(s.Service).Topics.Publish(topic, request)

	publish.Context(ctx)
	response, err := publish.Do()
	if err != nil {
		return nil, err
	}
	return response.MessageIds, nil
}

//New creates a service
func New(ctx context.Context, projectId string) (msgbus.Service, error) {
	srv, err := pubsub.NewService(ctx)
	if err != nil {
		return nil, err
	}
	return &service{
		projectID:projectId,
		Service: srv,
	}, nil
}
