package cmd

import (
	"context"
	"github.com/viant/afs/storage"
	"github.com/viant/afs/url"
	"github.com/viant/smirror/cmd/history"
	"github.com/viant/smirror/cmd/mirror"
	"github.com/viant/smirror/contract"
	"time"
)

func (s *service) emit(ctx context.Context, object storage.Object, req *mirror.Request, response *mirror.Response) error {
	if object.IsDir() {
		eventsHistory, err := history.FromURL(ctx, req.HistoryPathURL(object.URL()), s.fs)
		defer func() {
			eventsHistory.Persist(ctx, s.fs)
		}()
		objects, err := s.fs.List(ctx, object.URL())
		if err != nil {
			return err
		}
		for i := range objects {
			if url.Equals(object.URL(), objects[i].URL()) {
				continue
			}

			if ! object.IsDir()  && ! eventsHistory.Put(history.NewSource(object.URL(), object.ModTime())) {
				continue
			}
			if err := s.emit(ctx, objects[i], req, response); err != nil {
				return err
			}
		}
		return nil
	}

	request := &contract.Request{
		URL: object.URL(),
		Timestamp:time.Now(),
	}
	response.AddDataURLs(object.URL())
	response.IncrementPending(1)
	s.requestChan <- request
	return nil
}
