package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/viant/afsc/s3"
	"smirror/base"
	"smirror/cron"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered in f", r)
		}
	}()
	lambda.Start(handleRequest)
}



func handleRequest(ctx context.Context) (*cron.Response, error) {
	if base.IsLoggingEnabled() {
		s3.Logf = s3.StdoutLogger
	}
	service, err := cron.NewFromEnv(ctx, base.ConfigEnvKey)
	if err != nil {
		return nil, err
	}
	response := service.Tick(ctx)
	if data, err := json.Marshal(response); err == nil {
		fmt.Printf("%s\n", data)
	}
	return response, nil
}
