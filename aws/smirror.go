package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/viant/afsc/gs"
	_ "github.com/viant/afsc/s3"
	"runtime/debug"
	"smirror"
	"smirror/base"
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
	service, err := smirror.NewFromEnv(ctx, base.ConfigEnvKey)
	if err != nil {
		return err
	}

	for _, resource := range s3Event.Records {
		URL := resourceURL(resource)
		if base.IsLoggingEnabled() {
			fmt.Printf("triggered by  %v\n", URL)
		}
		response := service.Mirror(ctx, smirror.NewRequest(URL))
		if base.IsLoggingEnabled() {
			if data, err := json.Marshal(response); err == nil {
				fmt.Printf("%s\n", string(data))
			}
		}

	}
	return nil
}

//resourceURL returns resource URL
func resourceURL(resource events.S3EventRecord) string {
	return fmt.Sprintf("s3://%s/%s", resource.S3.Bucket.Name, resource.S3.Object.Key)
}
