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
func (s *service) Publish(ctx context.Context, request *msgbus.Request) (*msgbus.Response, error) {
	response := &msgbus.Response{MessageIDs: make([]string, 0)}
	return response, s.publish(ctx, request, response)
}

func (s *service) publish(ctx context.Context, request *msgbus.Request, response *msgbus.Response) error {
	publishRequest := &pubsub.PublishRequest{
		Messages: []*pubsub.PubsubMessage{
			{
				Data: base64.StdEncoding.EncodeToString(request.Data),
			},
		},
	}
	if len(request.Attributes) > 0 {
		publishRequest.Messages[0].Attributes = make(map[string]string)
		for k, v := range request.Attributes {
			publishRequest.Messages[0].Attributes[k] = fmt.Sprintf("%s", v)
		}
	}

	topic := s.topicInProject(request.Dest)
	publishCall := pubsub.NewProjectsService(s.Service).Topics.Publish(topic, publishRequest)

	publishCall.Context(ctx)
	callResponse, err := publishCall.Do()
	if err == nil {
		response.MessageIDs = callResponse.MessageIds
	}
	return err
}

//New creates a service
func New(ctx context.Context, projectId string) (msgbus.Service, error) {
	srv, err := pubsub.NewService(ctx)
	if err != nil {
		return nil, err
	}
	return &service{
		projectID: projectId,
		Service:   srv,
	}, nil
}
