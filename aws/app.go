package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/viant/toolbox"
	"runtime/debug"
	"smirror"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			fmt.Println("Recovered in f", r)
		}
	}()
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, s3Event events.S3Event) error {
	if len(s3Event.Records) == 0 {
		return nil
	}
	service, err := smirror.NewFromEnv(smirror.ConfigEnvKey)
	if err != nil {
		return err
	}
	if smirror.IsFnLoggingEnabled(smirror.LoggingEnvKey) {
		fmt.Printf("uses service %p, %v\n", service, err)
	}

	for _, resource := range s3Event.Records {
		URL := resourceURL(resource)
		if smirror.IsFnLoggingEnabled(smirror.LoggingEnvKey) {
			fmt.Printf("triggered by  %v\n", URL)
		}
		response := service.Mirror(smirror.NewRequest(URL))
		toolbox.Dump(response)

	}
	return nil
}

//resourceURL returns resource URL
func resourceURL(resource events.S3EventRecord) string {
	return fmt.Sprintf("s3://%s/%s", resource.S3.Bucket.Name, resource.S3.Object.Key)
}
