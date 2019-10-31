package config

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/matcher"
	"github.com/viant/afs/storage"
	"smirror/base"
	"sync/atomic"
	"time"
)

//Ruleset represents resources rules to check for changes to trigger storage event
type Ruleset struct {
	BaseURL      string
	CheckInMs    int
	Rules        []*Rule
	initialRules []*Rule
	inited       int32
	projectID    string
	meta         *base.Meta
}

//Init initialises resources
func (r *Ruleset) Init(ctx context.Context, fs afs.Service, projectID string) error {
	r.initRules()
	r.projectID = projectID
	r.meta = base.NewMeta(r.BaseURL, time.Duration(r.CheckInMs)*time.Millisecond)
	return r.loadAndInit(ctx, fs)
}

func (r *Ruleset) loadAndInit(ctx context.Context, fs afs.Service) (err error) {
	if err = r.loadAllResources(ctx, fs); err != nil {
		return err
	}
	for i := range r.Rules {
		r.Rules[i].Source.Init(r.projectID)
		r.Rules[i].Dest.Init(r.projectID)
	}
	return nil
}


func (r *Ruleset) ReloadIfNeeded(ctx context.Context, fs afs.Service) (bool, error) {
	changed, err := r.meta.HasChanged(ctx, fs)
	if err != nil || !changed {
		return changed, err
	}
	return true, r.loadAndInit(ctx, fs)
}



func (r *Ruleset) loadAllResources(ctx context.Context, fs afs.Service) error {
	if r.BaseURL == "" {
		return nil
	}
	r.Rules = r.initialRules
	exists, err := fs.Exists(ctx, r.BaseURL)
	if err != nil || !exists {
		return err
	}

	suffixMatcher, _ := matcher.NewBasic("", ".json", "", nil)
	routesObject, err := fs.List(ctx, r.BaseURL, suffixMatcher)
	if err != nil {
		return err
	}
	for _, object := range routesObject {
		if object.IsDir() {
			continue
		}
		if err = r.loadResources(ctx, fs, object); err != nil {
			return err
		}
	}
	return nil
}

func (r *Ruleset) loadResources(ctx context.Context, storage afs.Service, object storage.Object) error {
	reader, err := storage.Download(ctx, object)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()
	resources := make([]*Rule, 0)
	err = json.NewDecoder(reader).Decode(&resources)
	if err != nil {
		return errors.Wrapf(err, "failed to decode: %v", object.URL())
	}
	r.Rules = append(r.Rules, resources...)
	return err
}

func (r *Ruleset) initRules() {
	if atomic.CompareAndSwapInt32(&r.inited, 0, 1) {
		if len(r.Rules) > 0 {
			r.initialRules = r.Rules
		} else {
			r.initialRules = make([]*Rule, 0)
		}
	}
}
