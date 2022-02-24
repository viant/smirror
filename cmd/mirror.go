package cmd

import (
	"context"
	"github.com/pkg/errors"
	"github.com/viant/afs/storage"
	"github.com/viant/smirror/base"
	"github.com/viant/smirror/cmd/history"
	"github.com/viant/smirror/cmd/mirror"
	"github.com/viant/smirror/config"
	"github.com/viant/smirror/contract"
	"github.com/viant/smirror/shared"
	"sync"
	"sync/atomic"
	"time"
)

func (s *service) Mirror(ctx context.Context, request *mirror.Request) (*mirror.Response, error) {
	if request.RuleURL == "" {
		return nil, errors.New("ruleURL was empty")
	}
	request.Init(s.config)
	rule, err := s.loadRule(ctx, request.RuleURL)
	if err != nil {
		return nil, err
	}
	s.reportRule(rule)
	object, err := s.fs.Object(ctx, request.SourceURL)
	if err != nil {
		return nil, errors.Wrapf(err, "source location not found: %v", request.SourceURL)
	}
	response := mirror.NewResponse()
	ctx, cancel := context.WithCancel(ctx)
	go s.mirrorInTheBackground(ctx, cancel)

	waitGroup := &sync.WaitGroup{}
	go s.handleResponse(ctx, response)

	for atomic.LoadInt32(&s.stopped) == 0 {
		s.loadDatafiles(waitGroup, ctx, object, rule, request, response)
		if !request.Stream {
			break
		}

		if len(response.DataURLs) > 0 {
			shared.LogLn(response)
			response = mirror.NewResponse()
		} else {
			time.Sleep(time.Second)
		}
	}

	s.Stop()
	return response, err
}

func (s *service) loadDatafiles(waitGroup *sync.WaitGroup, ctx context.Context, object storage.Object, rule *config.Rule, request *mirror.Request, response *mirror.Response) {
	waitGroup.Add(1)
	go s.scanFiles(ctx, waitGroup, object, request, response)
	waitGroup.Wait()

	for atomic.LoadInt32(&s.stopped) == 0 && response.Pending() > 0 {
		time.Sleep(2 * time.Second)
		shared.LogProgress()
	}
	if len(response.Errors) == 0 {
		if err := s.updateHistory(ctx, response); err != nil {
			response.AddError(err)
		}
	}
}

func (s *service) scanFiles(ctx context.Context, waitGroup *sync.WaitGroup, object storage.Object, request *mirror.Request, response *mirror.Response) {
	defer waitGroup.Done()
	if err := s.emit(ctx, object, request, response); err != nil {
		response.AddError(err)
		s.Stop()
	}
	return
}

func (s *service) handleResponse(ctx context.Context, response *mirror.Response) {
	for {
		select {
		case <-s.stopChan:
			return
		case resp := <-s.responseChan:
			if resp.Error != "" {
				s.Stop()
				response.AddError(errors.New(resp.Error))
				return
			}

			if len(resp.MessageIDs) > 0 {
				response.AddMessageIDs(resp.MessageIDs)
			}
			if resp.BadRecords > 0 {
				atomic.AddUint64(&response.BadRecords, uint64(resp.BadRecords))
			}
			response.IncrementPending(-1)
			switch resp.Status {
			case base.StatusOK:
				response.AddMirrored(resp.TriggeredBy)
			case base.StatusError:
				response.AddFailed(resp.TriggeredBy)
			case base.StatusNoMatch:
				response.AddNoMatch(resp.TriggeredBy)
			}
		}
	}
}

func (s *service) mirrorInTheBackground(ctx context.Context, cancel context.CancelFunc) {
	for {
		select {
		case <-s.stopChan:
			cancel()
			return
		case req := <-s.requestChan:
			go func(req *contract.Request) {
				resp := s.mirrorService.Mirror(ctx, req)
				s.responseChan <- resp
			}(req)
		}
	}
}

func (s *service) updateHistory(ctx context.Context, response *mirror.Response) error {
	historyURLs := response.HistoryURLs()
	if len(historyURLs) == 0 {
		return nil
	}
	var index = make(map[string]bool)
	for _, URL := range response.DataURLs {
		index[URL] = true
	}
	for _, URL := range historyURLs {
		events, err := history.FromURL(ctx, URL, s.fs)
		if err != nil {
			return err
		}
		for i, event := range events.Events {
			if index[event.URL] {
				events.Events[i].Status = base.StatusOK
			}
		}
		err = events.Persist(ctx, s.fs)
		if err != nil {
			return err
		}
	}
	return nil
}
