package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"os"
	"smirror/aws/proxy"
)

const DestEnvKey = "DEST"

func handleMessages(ctx context.Context, sqsEvent events.SNSEvent) (err error) {
	dest := os.Getenv(DestEnvKey)
	if dest == "" {
		log.Printf("env.%v key was empty", DestEnvKey)
	}
	if len(sqsEvent.Records) == 0 {
		return err
	}
	proxier, err := proxy.Singleton()
	if err != nil {
		return err
	}
	for _, record := range sqsEvent.Records {
		response := proxier.Do(ctx, dest, []byte(record.SNS.Message))
		if data, err := json.Marshal(response); err == nil {
			fmt.Printf("%v\n", string(data))
		}
	}
	return err
}

func main() {
	lambda.Start(handleMessages)
}
