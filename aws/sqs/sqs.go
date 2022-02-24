package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/viant/smirror/event"
	"github.com/viant/smirror/proxy"
)

var config *proxy.Config

func handleMessages(ctx context.Context, sqsEvent events.SQSEvent) (err error) {
	if len(sqsEvent.Records) == 0 {
		return nil
	}
	if config == nil {
		config, err = proxy.NewConfig(ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to create config")
		}
	}
	proxier := proxy.Singleton(config)
	for _, record := range sqsEvent.Records {
		s3Event, err := event.NewS3EventFromJSON([]byte(record.Body))
		if err != nil {
			return errors.Wrapf(err, "unable unmarshal s3 event from %s", record.Body)
		}
		err = s3Event.Each(func(URL string) error {
			response := proxier.Proxy(ctx, &proxy.Request{
				Source: config.Source.CloneWithURL(URL),
				Dest:   &config.Dest,
				Move:   config.Move,
			})
			response.SourceURL = URL
			if data, err := json.Marshal(response); err == nil {
				fmt.Printf("%v\n", string(data))
			}
			if response.Error != "" {
				return errors.New(response.Error)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return err
}

func main() {
	lambda.Start(handleMessages)
}
