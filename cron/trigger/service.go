package trigger

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/option"
	"github.com/viant/afs/storage"
	"github.com/viant/afs/url"
	"smirror/base"
	"smirror/cron/config"
	"smirror/secret"
	"strings"
	"time"
)

//Service reresents trigger service
type Service interface {
	Trigger(ctx context.Context, resource *config.Rule, eventSource storage.Object) (map[string]string, error)
}

type service struct {
	*lambda.Lambda
	fs     afs.Service
	secret secret.Service
}

//Trigger triggers lambda execution
func (s *service) Trigger(ctx context.Context, resource *config.Rule, eventSource storage.Object) (map[string]string, error) {
	if url.Scheme(resource.Dest, "") == "" {
		return nil, s.notifyLambda(ctx, resource, eventSource)
	}
	var options = make([]storage.Option, 0)
	sourceOptions, err := s.secret.StorageOpts(ctx, &resource.Source)
	if err != nil {
		return nil, err
	}
	if len(sourceOptions) > 0 {
		options = append(options, option.NewSource(sourceOptions...))
		options = append(options, option.NewDest(option.NewAuth(true)))
	}
	sourceBucket := url.Host(eventSource.URL())
	destURL := url.Join(resource.Dest, sourceBucket)
	var transferred = make(map[string]string)
	return transferred, base.Trigger(ctx, s.fs, resource.Move, eventSource.URL(), destURL, transferred)
}

func (s *service) notifyLambda(ctx context.Context, resource *config.Rule, eventSource storage.Object) error {
	URL := eventSource.URL()
	bucket := url.Host(URL)
	URLPath := url.Path(URL)
	s3Event := events.S3Event{Records: make([]events.S3EventRecord, 0)}
	s3Event.Records = append(s3Event.Records, events.S3EventRecord{
		AWSRegion:   resource.Source.Region,
		EventTime:   time.Now(),
		EventSource: "s3",
		S3: events.S3Entity{
			Bucket: events.S3Bucket{
				Name: bucket,
			},
			Object: events.S3Object{
				Key:  strings.Trim(URLPath, "/"),
				Size: eventSource.Size(),
			},
		},
	})
	payload, err := json.Marshal(s3Event)
	if err != nil {
		return errors.Wrapf(err, "failed to decode s3 event for %v", eventSource.URL())
	}
	input := &lambda.InvokeInput{
		FunctionName:   &resource.Dest,
		Payload:        payload,
		InvocationType: aws.String(lambda.InvocationTypeEvent),
	}
	_, err = s.Invoke(input)
	return err
}

//New create trigger service
func New(fs afs.Service, secret secret.Service) (Service, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return &service{Lambda: lambda.New(sess), fs: fs, secret: secret}, nil
}
