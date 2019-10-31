package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/option"
	"github.com/viant/afs/storage"
	"github.com/viant/afs/url"
	"path"
	"smirror/base"
	"smirror/event"
	"smirror/secret"
)

const partSize = 32 * 1024 * 1024

//Service represents trigger service
type Service interface {
	Proxy(ctx context.Context, request *Request) (*Response)
}

type service struct {
	fs     afs.Service
	secret secret.Service
}

//Trigger triggers lambda execution
func (s *service) Proxy(ctx context.Context, request *Request) *Response {
	respose := NewResponse()
	if err := s.proxy(ctx, request, respose); err != nil {
		respose.Status = base.StatusError
		respose.Error = err.Error()
	}
	return respose
}

//Trigger triggers lambda execution
func (s *service) proxy(ctx context.Context, request *Request, response *Response) error {
	if err := request.Validate(); err != nil {
		return err
	}

	scheme := url.Scheme(request.Dest.URL, "")
	switch scheme {
	case base.LambdaScheme:
		return invokeLambda(ctx, request, response)
	case base.CloudFunctionScheme:
		return s.invokeCloudFunction(ctx, request, response)

	}
	var options = make([]storage.Option, 0)
	sourceOptions, err := s.secret.StorageOpts(ctx, request.Source)
	if err != nil {
		return err
	}
	if len(sourceOptions) == 0 {
		sourceOptions = make([]storage.Option, 0)
	}
	destOptions, err := s.secret.StorageOpts(ctx, request.Dest)
	if err != nil {
		return err
	}
	if len(destOptions) == 0 {
		destOptions = make([]storage.Option, 0)
	}
	if len(sourceOptions) > 0 {
		destOptions = append(destOptions, option.NewAuth(true))
		options = append(options, option.NewSource(sourceOptions...))
	}
	if len(destOptions) > 0 {
		options = append(options, option.NewDest(destOptions...))
	}
	if request.Stream {
		object, err := s.fs.Object(ctx, request.Source.URL, sourceOptions...)
		if err != nil {
			return errors.Wrapf(err, "source not found: %v", request.Source.URL)
		}
		sourceOptions = append(sourceOptions, option.NewStream(partSize, int(object.Size())))
		destOptions = append(destOptions, option.NewChecksum(true))
	}
	sourceBucket := url.Host(request.Source.URL)
	_, sourcePath := url.Base(request.Source.URL, "")
	destURL := url.Join(request.Dest.URL, path.Join(sourceBucket, sourcePath))
	transferred := response.Copied
	if request.Move {
		transferred = response.Moved
	}
	return s.propagate(ctx, request.Move, request.Source.URL, destURL, transferred, options...)
}

func invokeLambda(ctx context.Context, request *Request, response *Response) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	function := url.Host(request.Dest.URL)
	destService := lambda.New(sess)
	s3event := event.NewS3EventForURL(request.Source.URL)
	payload, err := json.Marshal(s3event)
	if err != nil {
		return err
	}
	response.AddInvoked(function, string(payload))
	_, err = destService.Invoke(&lambda.InvokeInput{
		FunctionName:   &function,
		Payload:        payload,
		InvocationType: aws.String(lambda.InvocationTypeEvent),
	})
	return err
}


//Note that invoking cloud function has restriction and should not be used on production
func (s *service) invokeCloudFunction(ctx context.Context, request *Request, response *Response) error {
	return fmt.Errorf("calling cloud function is not yet supported")
}

//propagate propagate source event with copy or move operation both are stream
func (s *service) propagate(ctx context.Context, isMove bool, sourceURL, destURL string, triggered map[string]string, options ...storage.Option) error {
	triggerFunc := s.fs.Copy
	if isMove {
		triggerFunc = s.fs.Move
	}
	triggered[sourceURL] = destURL

	err := triggerFunc(ctx, sourceURL, destURL, options...)
	if err != nil {
		if exists, e := s.fs.Exists(ctx, sourceURL); e == nil && !exists {
			err = nil
			triggered[sourceURL] = base.StatusNoFound
		}
	}
	return err
}

//New create trigger service
func New(fs afs.Service, secret secret.Service) Service {
	return &service{fs: fs, secret: secret}
}
