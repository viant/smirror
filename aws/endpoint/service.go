package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/viant/afs"
	"log"
	"os"
	"smirror"
	"smirror/base"
	"smirror/contract"
	"smirror/event"
	"sync"
	"time"
)

//Service represents subscriber service
type Service struct {
	config    *Config
	fs        afs.Service
	awsConfig *aws.Config
	session   *session.Session
	sqs       *sqs.SQS
	smirror.Service
	mux sync.Mutex
}

//Consume starts consumer
func (s *Service) Consume(ctx context.Context) error {
	for {
		err := s.consume(ctx)
		if err != nil {
			log.Printf("failed to consume: %v\n", err)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func (s *Service) consume(ctx context.Context) error {
	var URL string
	defer func() {
		r := recover()
		if r != nil {
			fmt.Printf("recover from panic: URL:%v, error: %v", URL, r)
		}
	}()
	queueURL, err := s.getQueueURL()
	if err != nil {
		log.Panicf("invalid queue: %v %v", s.config.Queue, err)
	}
	i := 0
	pullCount := s.config.BatchSize
	receivedInput := buildReceiveMessageInput(queueURL, pullCount, s.config.WaitTimeSec,  s.config.VisibilityTimeout,true)
	output, err := s.sqs.ReceiveMessage(receivedInput)
	if err != nil {
		return fmt.Errorf("failed to receive queue messages: %v, %w", queueURL, err)
	}
	i += len(output.Messages)
	deleteInput := &sqs.DeleteMessageBatchInput{
		Entries:  make([]*sqs.DeleteMessageBatchRequestEntry, 0),
		QueueUrl: aws.String(queueURL),
	}
	waitGroup := sync.WaitGroup{}
	for _, msg := range output.Messages {
		waitGroup.Add(1)
		go s.handleMessageInBackground(ctx, msg, deleteInput, &waitGroup)
	}
	waitGroup.Wait()
	if len(deleteInput.Entries) > 0 {
		_, err = s.sqs.DeleteMessageBatch(deleteInput)
	}
	if err != nil {
		return fmt.Errorf("failed to delete queue messages: %v, %w", queueURL, err)
	}
	return nil
}

func (s *Service) handleMessageInBackground(ctx context.Context, msg *sqs.Message, deleteInput *sqs.DeleteMessageBatchInput, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	ack, err := s.handleMessage(ctx, msg)
	if ack {
		s.mux.Lock()
		deleteInput.Entries = append(deleteInput.Entries, &sqs.DeleteMessageBatchRequestEntry{
			Id:            msg.MessageId,
			ReceiptHandle: msg.ReceiptHandle,
		})
		s.mux.Unlock()
	}
	if err != nil {
		fmt.Printf("failed to handle message: %v: (body: %s)%v\n", *msg.MessageId, *msg.Body, err)
	}
}

func (s *Service) handleMessage(ctx context.Context, msg *sqs.Message) (bool, error) {
	s3Event := event.S3Event{}
	if msg.Body == nil {
		return true, fmt.Errorf("message body was empty %v", *msg.MessageId)
	}
	if err := json.Unmarshal([]byte(*msg.Body), &s3Event); err != nil {
		return true, fmt.Errorf("failed to unmarshal s3Event: %s, due to %w", *msg.Body, err)
	}
	service, err := smirror.NewFromEnv(ctx, base.ConfigEnvKey)
	if err != nil {
		return false, err
	}
	if os.Getenv("DEBUG_MSG") != "" {
		fmt.Printf("%s\n", *msg.Body)
	}
	response := service.Mirror(ctx, contract.NewRequest(s3Event.URL()))
	output, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("failed marshal reported %v\n", response)
	}
	fmt.Printf("%s\n", output)
	return true, nil
}

func (s *Service) getQueueURL() (string, error) {
	result, err := s.sqs.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(s.config.Queue),
	})
	if err != nil {
		return "", fmt.Errorf("failed to lookup queue: %v, %w", s.config.Queue, err)
	}
	return *result.QueueUrl, nil
}

func buildReceiveMessageInput(queueURL string, pullCount int, waitTime int64, visibilityTimeout int64, includeAttr bool) *sqs.ReceiveMessageInput {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: aws.Int64(int64(pullCount)),
		WaitTimeSeconds:     aws.Int64(waitTime),
		VisibilityTimeout:   aws.Int64(visibilityTimeout),
	}
	if includeAttr {
		input.MessageAttributeNames = aws.StringSlice([]string{"All"})
		input.AttributeNames = aws.StringSlice([]string{"All"})
	}
	return input
}

//New creates a new subscriber
func New(config *Config, awsConfig *aws.Config, fs afs.Service) (*Service, error) {
	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	return &Service{
		config:    config,
		fs:        fs,
		sqs:       sqs.New(awsSession),
		session:   awsSession,
		awsConfig: awsConfig,
	}, nil
}
