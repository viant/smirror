package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/viant/smirror/base"
	"github.com/viant/smirror/cron"
	"github.com/viant/smirror/shared"
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
	service, err := cron.NewFromEnv(ctx, base.ConfigEnvKey)
	if err != nil {
		return nil, err
	}
	response := service.Tick(ctx)
	shared.LogLn(response)
	return response, nil
}
