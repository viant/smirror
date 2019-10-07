package mon

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/matcher"
	"github.com/viant/afs/option"
	"github.com/viant/afs/storage"
	"io/ioutil"
	"smirror/base"
	"time"
)

//Service represents monitoring service
type Service interface {
	//Check checks un process file and mirror errors
	Check(context.Context, *Request) *Response
}

type service struct {
	fs afs.Service
	*Config
}

//Check checks triggerBucket and error
func (s *service) Check(ctx context.Context, request *Request) *Response {
	response := NewResponse()
	err := s.check(ctx, request, response)
	if err != nil {
		response.Error = err.Error()
		response.Status = base.StatusError
	} else if response.UnprocessedCount > 0 {
		response.Status = base.StatusUnProcess
	} else if len(response.Errors) > 0 {
		response.Status = base.StatusError
		response.Error = response.Errors[0].Message
	}
	return response
}

func (s *service) check(ctx context.Context, request *Request, response *Response) (err error) {
	if err = request.Init(); err != nil {
		return err
	}
	if request.ErrorURL != "" {
		if err = s.checkErrors(ctx, request, response); err != nil {
			return err
		}
	}
	if request.ProcessedURL != "" {
		if err = s.checkProcessed(ctx, request, response); err != nil {
			return err
		}
	}
	return s.checkUnprocessed(ctx, request, response)
}

func (s *service) list(ctx context.Context, URL string, modifiedBefore, modifiedAfter *time.Time) ([]storage.Object, error) {
	timeMatcher := matcher.NewModification(modifiedBefore, modifiedAfter)
	recursive := option.NewRecursive(true)
	exists, _ := s.fs.Exists(ctx, URL)
	if !exists {
		return []storage.Object{}, nil
	}
	return s.fs.List(ctx, URL, timeMatcher, recursive)
}

func (s *service) checkErrors(ctx context.Context, request *Request, response *Response) error {
	objects, err := s.list(ctx, request.ErrorURL, nil, request.errorModifiedAfter)
	if err != nil {
		return errors.Wrapf(err, "failed to check errors: %v", request.ErrorURL)
	}
	for _, object := range objects {
		if object.IsDir() {
			continue
		}
		reader, err := s.fs.Download(ctx, object)
		if err != nil {
			return err
		}
		message, err := ioutil.ReadAll(reader)
		_ = reader.Close()
		if err != nil {
			return err
		}
		response.AddError(object, string(message))
	}
	response.ErrorCount = len(response.Errors)
	return nil
}

func (s *service) loadRoutes(ctx context.Context, URL string) (*Routes, error) {
	reader, err := s.fs.DownloadWithURL(ctx, URL)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
	}()
	result := &Routes{}
	return result, json.NewDecoder(reader).Decode(&result)
}

func (s *service) checkProcessed(ctx context.Context, request *Request, response *Response) error {
	routes, err := s.loadRoutes(ctx, request.ConfigURL)
	if err != nil {
		return errors.Wrapf(err, "failed to load routes: configf from URL :%v", request.ConfigURL)
	}
	if err := routes.Mirrors.Init(ctx, s.fs, ""); err != nil {
		return err
	}
	objects, err := s.list(ctx, request.ProcessedURL, nil, request.processedModifiedAfter)
	if err != nil {
		return errors.Wrapf(err, "failed to check processed: %v", request.ProcessedURL)
	}
	for _, object := range objects {
		if object.IsDir() {
			continue
		}
		route := routes.Mirrors.HasMatch(object.URL())
		response.AddProcessed(route, object)
	}
	return nil
}

func (s *service) checkUnprocessed(ctx context.Context, request *Request, response *Response) error {
	routes, err := s.loadRoutes(ctx, request.ConfigURL)
	if err != nil {
		return errors.Wrapf(err, "failed to load routes: %v", request.ConfigURL)
	}
	if err := routes.Mirrors.Init(ctx, s.fs, ""); err != nil {
		return err
	}
	objects, err := s.list(ctx, request.TriggerURL, request.unprocessedModifiedBefore, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to check unprocessed: %v", request.TriggerURL)
	}
	now := time.Now()
	for _, object := range objects {
		if object.IsDir() {
			continue
		}
		route := routes.Mirrors.HasMatch(object.URL())
		response.AddUnprocessed(now, route, object)
	}
	return nil
}

//New creates new service
func New(config *Config) Service {
	config.Init()
	return &service{
		fs:     afs.New(),
		Config: config,
	}
}
