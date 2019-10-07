package sqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/viant/toolbox"
	"smirror/msgbus"
)

type service struct {
	session *session.Session
	sqs     *sqs.SQS
}

//Publish publisher data to sqs
func (s *service) Publish(ctx context.Context, request *msgbus.Request) (*msgbus.Response, error) {
	response := &msgbus.Response{MessageIDs: make([]string, 0)}
	return response, s.publish(ctx, request, response)
}

//Publish publishes data to sqs
func (s *service) publish(ctx context.Context, request *msgbus.Request, response *msgbus.Response) error {
	queueURL, err := s.getQueueURL(request.Dest)
	if err != nil {
		return err
	}
	input := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(1),
		QueueUrl:     &queueURL,
	}
	if len(request.Attributes) > 0 {
		input.MessageAttributes = make(map[string]*sqs.MessageAttributeValue)
		putSqsMessageAttributes(request.Attributes, input.MessageAttributes)
	}

	var body = toolbox.AsString(request.Data)
	input.MessageBody = aws.String(body)
	result, err := s.sqs.SendMessage(input)
	if err == nil {
		response.MessageIDs = []string{*result.MessageId}
	}
	return err
}

func putSqsMessageAttributes(attributes map[string]interface{}, target map[string]*sqs.MessageAttributeValue) {
	for k, v := range attributes {
		if v == nil {
			continue
		}
		dataType := getAttributeDataType(v)
		target[k] = &sqs.MessageAttributeValue{
			DataType:    &dataType,
			StringValue: aws.String(toolbox.AsString(v)),
		}
	}
}

func getAttributeDataType(value interface{}) string {
	dataType := "String"
	if toolbox.IsInt(value) || toolbox.IsFloat(value) {
		dataType = "Number"
	}
	return dataType
}

func (c *service) getQueueURL(queueName string) (string, error) {
	result, err := c.sqs.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to lookup queue URL %v", queueName)
	}
	return *result.QueueUrl, nil
}

//New creates a sqs service
func New(ctx context.Context) (msgbus.Service, error) {
	sess := session.New()
	return &service{
		session: sess,
		sqs:     sqs.New(sess),
	}, nil
}
