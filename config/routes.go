package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/matcher"
	"github.com/viant/afs/storage"
	"smirror/base"
	"strings"
	"sync/atomic"
	"time"
)

//Ruleset represents route slice
type Ruleset struct {
	BaseURL      string
	CheckInMs    int
	Rules        []*Rule
	meta         *base.Meta
	initialRules []*Rule
	inited       int32
}

//Match returns the first match route
func (r Ruleset) Match(URL string) (matched []*Rule) {
	for i := range r.Rules {
		if r.Rules[i].HasMatch(URL) {
			matched = append(matched, r.Rules[i])
		}
	}
	return matched
}

func (r Ruleset) Validate() error {
	if len(r.Rules) == 0 {
		return nil
	}
	for i := range r.Rules {
		if err := r.Rules[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}

//Init initialises resources
func (r *Ruleset) Init(ctx context.Context, fs afs.Service, projectID string) error {
	if err := r.initRules(); err != nil {
		return err
	}
	r.meta = base.NewMeta(r.BaseURL, time.Duration(r.CheckInMs)*time.Millisecond)
	return r.load(ctx, fs)
}

func (r *Ruleset) load(ctx context.Context, fs afs.Service) (err error) {
	if err = r.loadAllResources(ctx, fs); err != nil {
		return err
	}
	return nil
}

func (r *Ruleset) ReloadIfNeeded(ctx context.Context, fs afs.Service) (bool, error) {
	changed, err := r.meta.HasChanged(ctx, fs)
	if err != nil || !changed {
		return changed, err
	}
	return true, r.load(ctx, fs)
}

func (c *Ruleset) loadAllResources(ctx context.Context, fs afs.Service) error {
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
		if err = c.loadResources(ctx, fs, object); err != nil {
			//Report error, let the other rules work fine
			fmt.Println(err)
		}
	}
	return nil
}

func (c *Ruleset) loadResources(ctx context.Context, storage afs.Service, object storage.Object) error {
	reader, err := storage.Download(ctx, object)
	defer func() {
		_ = reader.Close()
	}()
	routes := make([]*Rule, 0)
	err = json.NewDecoder(reader).Decode(&routes)
	if err != nil {
		return errors.Wrapf(err, "failed to decode: %v", object.URL())
	}
	transientRoutes := Ruleset{Rules: routes}
	if err := transientRoutes.Validate(); err != nil {
		return errors.Wrapf(err, "invalid rule: %v", object.URL())
	}

	for i := range routes {
		routes[i].Info.URL = object.URL()
		if routes[i].Info.Workflow == "" {

			name := object.Name()
			if strings.HasSuffix(name, ".json") {
				name = string(name[:len(name)-5])
			}
			routes[i].Info.Workflow = name
		}
		c.Rules = append(c.Rules, routes[i])

	}

	return nil
}

func (r *Ruleset) initRules() error {
	if atomic.CompareAndSwapInt32(&r.inited, 0, 1) {
		if len(r.Rules) > 0 {
			if err := r.Validate(); err != nil {
				return err
			}
			r.initialRules = r.Rules
		} else {
			r.initialRules = make([]*Rule, 0)
		}
	}
	return nil
}
