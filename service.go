package smirror

import (
	"context"
	"fmt"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/storage"
	"github.com/viant/toolbox/storage/gs"
	"github.com/viant/toolbox/storage/s3"
	"io"
	"smirror/job"
	"smirror/secret"
	"sync"
	"sync/atomic"
	"time"
)

//Service represents a mirror service
type Service interface {
	//Mirror copies/split source to matched destination
	Mirror(request *Request) *Response
}

type service struct {
	config *Config
}

func (s *service) Mirror(request *Request) *Response {
	response := NewResponse()
	if err := s.mirror(request, response); err != nil {
		response.Status = StatusError
		response.Error = err.Error()
	}
	response.TimeTakenMs = int(time.Now().Sub(response.startTime) / time.Millisecond)
	return response

}

func (s *service) mirror(request *Request, response *Response) error {
	route := s.config.Routes.HasMatch(request.URL)
	if route == nil {
		response.Status = StatusNoMatch
		return nil
	}
	storageService, err := storage.NewServiceForURL(request.URL, "")
	if err != nil {
		return err
	}
	if route.Split != nil {
		err = s.mirrorChunkeddAsset(route, storageService, request, response)
	} else {
		err = s.mirrorAsset(route, request.URL, storageService, response)
	}
	context := job.NewContext(err, storageService, request.URL)
	if e := route.OnCompletion.Run(context); e != nil && err == nil {
		err = e
	}
	return err
}

func (s *service) mirrorAsset(route *Route, URL string, storageService storage.Service, response *Response) error {
	reader, err := storageService.DownloadWithURL(URL)
	sourceCompression := NewCompressionForURL(URL)
	destCompression := route.Compression
	if sourceCompression.Equals(destCompression) {
		sourceCompression = nil
		destCompression = nil
	}
	if err == nil {
		reader, err = NewReader(reader, sourceCompression)
		if err != nil {
			return err
		}
	}
	defer func() {
		if reader != nil {
			_ = reader.Close()
		}
	}()

	destName := route.Name(URL)
	destURL := toolbox.URLPathJoin(route.DestURL, destName)
	dataCopy := &Copy{Reader: reader, Dest: NewDatafile(destURL, destCompression)}
	return s.copy(dataCopy, response)
}

func (s *service) mirrorChunkeddAsset(route *Route, storageService storage.Service, request *Request, response *Response) error {
	reader, err := storageService.DownloadWithURL(request.URL)
	if err == nil {
		reader, err = NewReader(reader, NewCompressionForURL(request.URL))
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
	err = toolbox.SplitTextStream(reader, s.chunkWriter(request.URL, route, &counter, waitGroup, response), route.Split.MaxLines)
	if err == nil {
		waitGroup.Wait()
	}

	context := job.NewContext(err, storageService, request.URL)
	if e := route.OnCompletion.Run(context); e != nil {
		err = e
	}
	return err
}

func (s *service) copy(copy *Copy, response *Response) error {
	service, err := storage.NewServiceForURL(copy.Dest.URL, "")
	if err != nil {
		return err
	}
	reader, err := copy.GetReader()
	if err != nil {
		return err
	}
	if err = service.Upload(copy.Dest.URL, reader); err != nil {
		return err
	}
	response.AddURL(copy.Dest.URL)

	return nil
}

func (s *service) chunkWriter(URL string, route *Route, counter *int32, waitGroup *sync.WaitGroup, response *Response) func() io.WriteCloser {
	return func() io.WriteCloser {
		splitCount := atomic.AddInt32(counter, 1)
		destName := route.Split.Name(route, URL, splitCount)
		destURL := toolbox.URLPathJoin(route.DestURL, destName)

		return NewWriter(route, func(writer *Writer) error {
			waitGroup.Add(1)
			defer waitGroup.Done()
			dataCopy := &Copy{
				Reader: writer.Reader,
				Dest:   NewDatafile(destURL, nil),
			}
			return s.copy(dataCopy, response)
		})
	}
}

func (s *service) initSecrets() error {
	if len(s.config.Secrets) == 0 {
		return nil
	}

	for i, config := range s.config.Secrets {
		credConfig, err := secret.New(context.Background(), s.config.Secrets[i])
		if err != nil {
			return err
		}
		switch config.TargetScheme {
		case "gs":
			gs.SetProvider(credConfig)
		case "s3":
			s3.SetProvider(credConfig)
		default:
			return fmt.Errorf("unsupported target scheme: %v", config.TargetScheme)
		}
	}
	return nil
}

//New creates a new mirror service
func New(config *Config) (Service, error) {
	result := &service{config: config}
	err := result.initSecrets()
	if err != nil {
		return nil, err
	}
	return result, err
}
