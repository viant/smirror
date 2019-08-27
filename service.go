package smirror

import (
	"context"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	"github.com/viant/toolbox"
	"io"
	"smirror/config"
	"smirror/job"
	"smirror/secret"
	"sync"
	"sync/atomic"
	"time"
)

//Service represents a mirror service
type Service interface {
	//Mirror copies/split source to matched destination
	Mirror(ctx context.Context, request *Request) *Response
}

type service struct {
	config *Config
	afs.Service
	secret secret.Service
}

func (s *service) Mirror(ctx context.Context, request *Request) *Response {
	response := NewResponse()
	if err := s.mirror(ctx, request, response); err != nil {
		response.Status = StatusError
		response.Error = err.Error()
	}
	response.TimeTakenMs = int(time.Now().Sub(response.startTime) / time.Millisecond)
	return response

}

func (s *service) mirror(ctx context.Context, request *Request, response *Response) (err error) {
	route := s.config.Routes.HasMatch(request.URL)
	if route == nil {
		response.Status = StatusNoMatch
		return nil
	}
	if route.Split != nil {
		err = s.mirrorChunkedAsset(ctx, route, request, response)
	} else {
		err = s.mirrorAsset(ctx, route, request.URL, response)
	}
	jobContent := job.NewContext(ctx, err, request.URL)
	if e := route.OnCompletion.Run(jobContent, s.Service); e != nil && err == nil {
		err = e
	}
	return err
}

func (s *service) mirrorAsset(ctx context.Context, route *config.Route, URL string, response *Response) error {
	options, err := s.secret.StorageOpts(ctx, route.Source)
	if err != nil {
		return err
	}
	reader, err := s.Service.DownloadWithURL(ctx, URL, options...)
	if err != nil {
		return err
	}

	sourceCompression := config.NewCompressionForURL(URL)
	destCompression := route.Compression
	if sourceCompression.Equals(destCompression) {
		sourceCompression = nil
		destCompression = nil
	}
	reader, err = NewReader(reader, sourceCompression)
	if err != nil {
		return err
	}
	defer func() {
		if reader != nil {
			_ = reader.Close()
		}
	}()
	destName := route.Name(URL)
	destURL := url.Join(route.Dest.URL, destName)

	if reader == nil {
		return fmt.Errorf("reader was empty")
	}

	dataCopy := &Copy{
		Resource: &route.Dest,
		Reader:   reader,
		Dest:     NewDatafile(destURL, destCompression)}
	return s.copy(ctx, dataCopy, response)
}

func (s *service) mirrorChunkedAsset(ctx context.Context, route *config.Route, request *Request, response *Response) error {
	options, err := s.secret.StorageOpts(ctx, route.Source)
	if err != nil {
		return err
	}
	reader, err := s.Service.DownloadWithURL(ctx, request.URL, options...)
	if err == nil {
		reader, err = NewReader(reader, config.NewCompressionForURL(request.URL))
		if err != nil {
			return err
		}
	}
	defer func() {
		if reader != nil {
			_ = reader.Close()
		}
	}()
	counter := int32(0)
	waitGroup := &sync.WaitGroup{}
	err = toolbox.SplitTextStream(reader, s.chunkWriter(ctx, request.URL, route, &counter, waitGroup, response), route.Split.MaxLines)
	if err == nil {
		waitGroup.Wait()
	}
	jobContent := job.NewContext(ctx, err, request.URL)
	if e := route.OnCompletion.Run(jobContent, s.Service); e != nil {
		err = e
	}
	return err
}

func (s *service) copy(ctx context.Context, copy *Copy, response *Response) error {

	reader, err := copy.GetReader()
	if err != nil {
		return err
	}
	options, err := s.secret.StorageOpts(ctx, copy.Resource)
	if err != nil {
		return err
	}
	if err = s.Service.Upload(ctx, copy.Dest.URL, file.DefaultFileOsMode, reader, options...); err != nil {
		return err
	}
	response.AddURL(copy.Dest.URL)

	return nil
}

func (s *service) chunkWriter(ctx context.Context, URL string, route *config.Route, counter *int32, waitGroup *sync.WaitGroup, response *Response) func() io.WriteCloser {
	return func() io.WriteCloser {
		splitCount := atomic.AddInt32(counter, 1)
		destName := route.Split.Name(route, URL, splitCount)
		destURL := toolbox.URLPathJoin(route.Dest.URL, destName)
		return NewWriter(route, func(writer *Writer) error {
			if writer.Reader == nil {
				return fmt.Errorf("writer reader was empty")
			}
			waitGroup.Add(1)
			defer waitGroup.Done()
			dataCopy := &Copy{
				Resource: &route.Dest,
				Reader:   writer.Reader,
				Dest:     NewDatafile(destURL, nil),
			}
			return s.copy(ctx, dataCopy, response)
		})
	}
}

//New creates a new mirror service
func New(ctx context.Context, config *Config) (Service, error) {
	err := config.Init()
	if err != nil {
		return nil, err
	}
	result := &service{config: config,
		Service: afs.New(),
		secret:  secret.New(config.SourceScheme)}
	if resources := config.Resources(); len(resources) > 0 {
		err = result.secret.Init(ctx, result.Service, resources)
	}
	return result, err
}
