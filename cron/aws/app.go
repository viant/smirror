package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"smirror"
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

func handleRequest(ctx context.Context) error {
	service, err := cron.NewFromEnv(ctx, smirror.ConfigEnvKey)
	if err != nil {
		return err
	}
	if base.IsLoggingEnabled() {
		fmt.Printf("uses service %T(%p), err: %v\n", service, service, err)
	}
	err = service.Tick(ctx)
	if err != nil {
		log.Print(err)
	}
	return nil
}
