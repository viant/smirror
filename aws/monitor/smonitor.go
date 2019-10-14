package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"runtime/debug"
	"smirror/base"
	"smirror/mon"
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

func handleRequest(ctx context.Context, request *mon.Request) (*mon.Response, error) {
	service, err := mon.NewFromEnv(base.ConfigEnvKey)
	if err != nil {
		return nil, err
	}
	return service.Check(ctx, request), nil
}
