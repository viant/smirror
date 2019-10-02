package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/matcher"
	"github.com/viant/afs/storage"
	"github.com/viant/toolbox"
	"smirror/base"
	"time"
)

//Resources represents resources rules to check for changes to trigger storage event
type Resources struct {
	BaseURL       string
	CheckInMs int
	Rules         []*Resource
	initialRules  []*Resource
	projectID     string
	meta          *base.Meta
}

//Init initialises resources
func (r *Resources) Init(ctx context.Context, fs afs.Service, projectID string) error {
	r.initRules()
	r.projectID = projectID
	r.meta = base.NewMeta(r.BaseURL, time.Duration(r.CheckInMs)*time.Millisecond)
	return r.loadAndInit(ctx, fs)
}

func (r *Resources) loadAndInit(ctx context.Context, fs afs.Service) (err error) {
	if err = r.loadAllResources(ctx, fs); err != nil {
		return err
	}
	for i := range r.Rules {
		r.Rules[i].Init(r.projectID)
	}
	toolbox.Dump(r.Rules)
	fmt.Printf("\n")
	return nil
}

func (r *Resources) ReloadIfNeeded(ctx context.Context, fs afs.Service) (bool, error) {
	changed, err := r.meta.HasChanged(ctx, fs)
	if err != nil || ! changed {
		return changed, err
	}
	if base.IsLoggingEnabled() {
		fmt.Printf("reloading rules\n")
	}
	return true, r.loadAndInit(ctx, fs)
}

func (c *Resources) loadAllResources(ctx context.Context, fs afs.Service) error {
	if c.BaseURL == "" {
		return nil
	}
	c.Rules = c.initialRules
	exists, err := fs.Exists(ctx, c.BaseURL)
	if err != nil || !exists {
		return err
	}
	suffixMatcher, _ := matcher.NewBasic("", ".json", "", nil)
	routesObject, err := fs.List(ctx, c.BaseURL, suffixMatcher)
	if err != nil {
		return err
	}
	for _, object := range routesObject {
		if object.IsDir() {
			continue
		}
		fmt.Printf("downloading: %v\n", object.URL())
		if err = c.loadResources(ctx, fs, object); err != nil {
			return err
		}
	}
	return nil
}

func (c *Resources) loadResources(ctx context.Context, storage afs.Service, object storage.Object) error {
	reader, err := storage.Download(ctx, object)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()
	resources := make([]*Resource, 0)
	err = json.NewDecoder(reader).Decode(&resources);
	if err != nil {
		return errors.Wrapf(err, "failed to decode: %v", object.URL())
	}
	return err
}

func (r *Resources) initRules() {
	if len(r.initialRules) == 0 {
		if len(r.Rules) > 0 {
			r.initialRules = r.Rules
		} else {
			r.initialRules = make([]*Resource, 0)
		}
	}
}
