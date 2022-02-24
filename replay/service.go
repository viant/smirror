package replay

import (
	"context"
	"github.com/viant/smirror/base"
	"github.com/viant/afs"
	"github.com/viant/afs/matcher"
	"github.com/viant/afs/option"
	"github.com/viant/afs/storage"
	"time"
)

const replayExtension = ".replay"

//Service represents replay service
type Service interface {
	Replay(context.Context, *Request) *Response
}

type service struct {
	fs afs.Service
}

func (s *service) Replay(ctx context.Context, request *Request) *Response {
	response := &Response{
		Replayed: make([]string, 0),
		Status:   base.StatusOK,
	}
	err := s.replay(ctx, request, response)
	if err != nil {
		response.Status = base.StatusError
		response.Error = err.Error()
	}
	return response
}

func (s *service) replay(ctx context.Context, request *Request, response *Response) error {
	err := request.Init()
	if err == nil {
		err = request.Validate()
	}
	if err != nil {
		return err
	}
	objects, err := s.list(ctx, request.TriggerURL, request.unprocessedModifiedBefore)
	replayer := base.NewReplayer(s.fs)
	replayer.Run(ctx, 10)
	for i := range objects {
		if objects[i].IsDir() {
			continue
		}
		replayer.Schedule(objects[i].URL())
	}
	return replayer.Wait()
}

func (s *service) list(ctx context.Context, URL string, modifiedBefore *time.Time) ([]storage.Object, error) {
	timeMatcher := matcher.NewModification(modifiedBefore, nil)
	recursive := option.NewRecursive(true)
	exists, _ := s.fs.Exists(ctx, URL)
	if !exists {
		return []storage.Object{}, nil
	}
	return s.fs.List(ctx, URL, timeMatcher, recursive)
}

//New creates new replay service
func New() Service {
	return &service{
		fs: afs.New(),
	}
}
