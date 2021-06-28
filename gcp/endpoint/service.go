package endpoint

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/viant/afs"
	"golang.org/x/oauth2/google"
	"log"
	"os"
	"smirror"
	"smirror/base"
	"smirror/contract"
	"smirror/event"
	"strings"
	"sync"
	"time"
)

//Service represents client service
type Service struct {
	config *Config
	fs     afs.Service
	client *pubsub.Client
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

	subscription, err := s.getSubscription()
	if err != nil {
		return fmt.Errorf("failed to get subscription: %v, %w", s.config.Subscription, err)
	}
	subscription.ReceiveSettings.MaxOutstandingMessages = s.config.BatchSize
	subscription.ReceiveSettings.NumGoroutines = s.config.BatchSize
	subscription.ReceiveSettings.MaxExtension  = time.Duration(s.config.VisibilityTimeout) * time.Millisecond
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mutex := &sync.Mutex{}
	err = subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		mutex.Lock()
		defer mutex.Unlock()
		ok, err := s.handleMessage(ctx, msg)
		if err != nil {
			log.Printf("failed to handle message: %v", err)
		}
		if ok {
			msg.Ack()
		} else {
			msg.Nack()
		}

	})
	return nil
}

func (s *Service) handleMessage(ctx context.Context, msg *pubsub.Message) (bool, error) {
	gcsEvent := event.StorageEvent{}
	data := msg.Data
	if data == nil {
		return true, fmt.Errorf("message body was empty %+v", *msg)
	}

	if decoded, err := base64.StdEncoding.DecodeString(string(data)); err == nil {
		data = decoded
	}

	if err := json.Unmarshal([]byte(data), &gcsEvent); err != nil {
		return true, fmt.Errorf("failed to unmarshal gcsEvent: %s, due to %w", data, err)
	}
	service, err := smirror.NewFromEnv(ctx, base.ConfigEnvKey)
	if err != nil {
		return false, err
	}
	if os.Getenv("DEBUG_MSG") != "" {
		fmt.Printf("%s\n", data)
	}
	response := service.Mirror(ctx, contract.NewRequest(gcsEvent.URL()))
	output, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("failed marshal reported %v\n", response)
	}
	fmt.Printf("%s\n", output)
	return true, nil
}

func (s *Service) Close() error {
	return s.client.Close()
}

func (s *Service) getSubscription() (*pubsub.Subscription, error) {
	if s.config.Subscription == "" {
		return nil, fmt.Errorf("subscription was empty")
	}
	if !strings.Contains(s.config.Subscription, "/") {

	}
	if s.config.ProjectID == "" {
		return s.client.Subscription(s.config.Subscription), nil
	}
	return s.client.SubscriptionInProject(s.config.Subscription, s.config.ProjectID), nil
}

//New creates a new client
func New(config *Config, fs afs.Service) (*Service, error) {
	if config.ProjectID == "" {
		if credentials, err := google.FindDefaultCredentials(context.Background()); err == nil {
			config.ProjectID = credentials.ProjectID
		}
	}
	client, err := pubsub.NewClient(context.Background(), config.ProjectID)
	if err != nil {
		return nil, err
	}
	return &Service{
		config: config,
		client: client,
		fs:     fs,
	}, nil
}
