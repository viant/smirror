package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	flambda "github.com/aws/aws-sdk-go/service/lambda"
	"log"
	"os"
)

const DestEnvKey = "DEST"

func handleMessages(ctx context.Context, sqsEvent events.SQSEvent) (err error) {

	dest := os.Getenv(DestEnvKey)
	if dest == "" {
		log.Print("env.%v key was empty", DestEnvKey)
	}
	if len(sqsEvent.Records) == 0 {
		return err
	}
	for _, record := range sqsEvent.Records {
		if err = notify(dest, []byte(record.Body)); err != nil {
			return err
		}
	}
	return err
}

func main() {
	lambda.Start(handleMessages)
}

func notify(destination string, payload []byte) error {
	if destination == "" {
		log.Println("dest is empty, ingoring event: %s\n", payload)
		return nil
	}
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	service := flambda.New(sess)
	input := &flambda.InvokeInput{
		FunctionName:   &destination,
		Payload:        payload,
		InvocationType: aws.String(flambda.InvocationTypeEvent),
	}
	_, err = service.Invoke(input)
	return err
}
