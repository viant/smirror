package meta

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/storage"
	"time"
)

//Service represents meta service to managed process and pending resources
type Service interface {
	//PendingResources filters pending resources for supplied candidate with processed resources
	PendingResources(ctx context.Context, candidates []storage.Object) ([]storage.Object, error)

	//AddProcessed add processed resources
	AddProcessed(ctx context.Context, processed []storage.Object) error
}

type service struct {
	metaURL       string
	pruneDuration time.Duration
	afs.Service
}

func (s *service) loadState(ctx context.Context) (*State, error) {
	state := &State{}
	has, _ := s.Exists(ctx, s.metaURL)
	if !has {
		return state, nil
	}
	reader, err := s.OpenURL(ctx, s.metaURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load meta file: %v", s.metaURL)
	}
	err = json.NewDecoder(reader).Decode(state)
	return state, err
}

func (s *service) storeState(ctx context.Context, state *State) error {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(state)
	if err != nil {
		err = errors.Wrapf(err, "failed to encode meta: %v with %v", s.metaURL, state)
	}
	err = s.Upload(ctx, s.metaURL, 0644, buffer)
	if err != nil {
		err = errors.Wrapf(err, "failed to upload meta: %v", s.metaURL)
	}
	return err
}

//PendingResources filters pending resources for supplied candidate with processed resources
func (s *service) PendingResources(ctx context.Context, candidates []storage.Object) ([]storage.Object, error) {
	state, err := s.loadState(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load meta state")
	}
	var result = make([]storage.Object, 0)
	processed := state.ProcessMap()
	for i, candidate := range candidates {
		modified, has := processed[candidate.URL()]
		if !has || !modified.Equal(candidate.ModTime()) {
			result = append(result, candidates[i])
		}
	}
	return result, nil
}

//AddProcessed adds to processed resources
func (s *service) AddProcessed(ctx context.Context, processed []storage.Object) error {
	state, err := s.loadState(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to load meta state")
	}
	state.Add(processed...)
	state.Prune(time.Now(), s.pruneDuration)
	return s.storeState(ctx, state)
}

//New creates a new service
func New(metaURL string, pruneDuration time.Duration, fs afs.Service) Service {
	return &service{
		metaURL:       metaURL,
		pruneDuration: pruneDuration,
		Service:       fs,
	}
}
