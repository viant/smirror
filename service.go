package mirror

import (
	"smirror/job"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/storage"
	"io"
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
		err = s.mirrorSplittedAsset(route, storageService, request, response)
	} else {
		err = s.mirrorAsset(route, request.URL, storageService, response)
	}
	context := job.NewContext(err, storageService, request.URL)
	if e := route.OnCompletion.Run(context); e != nil && err == nil {
		err = e
	}
	return err
}


func (s *service) mirrorAsset(route *Route, URL string,  storageService storage.Service, response *Response) error {
	reader, err := storageService.DownloadWithURL(URL)
	//TODO optimze copy if dest uses the same compression, at the moment we decompress and comress it again
	if err == nil {
		reader, err = NewReader(reader, NewCompressionForURL(URL))
	}
	defer func() { _ = reader.Close() }()
	destName := route.Name(URL)
	destURL := toolbox.URLPathJoin(route.DestURL, destName)
	dataCopy := &Copy{Reader: reader, Dest: NewDatafile(destURL, route.Compression)}
	return s.copy(dataCopy, response)
}

func (s *service) mirrorSplittedAsset(route *Route, storageService storage.Service, request *Request, response *Response) error {
	reader, err := storageService.DownloadWithURL(request.URL)
	if err == nil {
		reader, err = NewReader(reader, NewCompressionForURL(request.URL))
	}
	defer func() { _ = reader.Close() }()
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

//New creates a new mirror service
func New(config *Config) Service {
	return &service{config:config}
}