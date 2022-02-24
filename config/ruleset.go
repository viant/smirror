package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/option"
	"github.com/viant/afs/storage"
	"github.com/viant/toolbox"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"github.com/viant/smirror/base"
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
func (r Ruleset) Rule(URL string) *Rule {
	for i := range r.Rules {
		if r.Rules[i].Info.URL == URL {
			return r.Rules[i]
		}
	}
	return nil
}


//Match returns the first match route
func (r Ruleset) Match(URL string) (matched []*Rule) {
	ruleURL := "."
	for i := range r.Rules {
		if r.Rules[i].HasMatch(URL) {
			if ruleURL == r.Rules[i].Info.URL {
				continue
			}
			ruleURL = r.Rules[i].Info.URL
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

func (r Ruleset) Init(ctx context.Context, fs afs.Service) error {
	if len(r.Rules) == 0 {
		return nil
	}
	for i := range r.Rules {
		if err := r.Rules[i].Init(ctx, fs); err != nil {
			return err
		}
	}
	return nil
}

//Load initialises resources
func (r *Ruleset) Load(ctx context.Context, fs afs.Service) error {
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
	fs.Delete(ctx,"s3://viant-dataflow-config/StorageMirror/_.cache")
	routesObject, err := fs.List(ctx, c.BaseURL, option.NewRecursive(true))
	if err != nil {
		return err
	}
	for _, object := range routesObject {
		if object.IsDir()  || ! (path.Ext(object.Name()) == ".json" || path.Ext(object.Name()) == ".yaml") {
			continue
		}

		if err = c.loadResources(ctx, fs, object); err != nil {
			//Report error, let the other rules work fine
			fmt.Println(err)
		}
	}
	return nil
}

func (c *Ruleset) loadResources(ctx context.Context, fs afs.Service, object storage.Object) error {
	reader, err := fs.Open(ctx, object)
	if err != nil {
		return fmt.Errorf("failed to open: %v, %w", object.URL(), err)
	}
	defer func() {
		if reader == nil {
			return
		}
		_ = reader.Close()
	}()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	rules, err := loadRules(data, path.Ext(object.Name()))
	if err != nil {
		return errors.Wrapf(err, "failed to load rules: %v", object.URL())
	}
	transientRoutes := Ruleset{Rules: rules}
	transientRoutes.Rules[0].Info.URL = object.URL()
	if err := transientRoutes.Init(ctx, fs); err != nil {
		return errors.Wrapf(err, "invalid rule: %v", object.URL())
	}
	if err := transientRoutes.Validate(); err != nil {
		return errors.Wrapf(err, "invalid rule: %v", object.URL())
	}
	for i := range rules {
		rules[i].Info.URL = object.URL()
		if rules[i].Info.Workflow == "" {
			name := object.Name()
			if strings.HasSuffix(name, ".json") {
				name = string(name[:len(name)-5])
			}
			rules[i].Info.Workflow = name
		}
		c.Rules = append(c.Rules, rules[i])
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

func loadRules(data []byte, ext string) ([]*Rule, error) {
	if ext == "" {
		return nil, nil
	}
	var rules = make([]*Rule, 0)
	switch ext {
	case base.YAMLExt:
		ruleMap := map[string]interface{}{}
		if err := yaml.Unmarshal(data, &ruleMap); err != nil {
			rulesMap := []map[string]interface{}{}
			err = json.Unmarshal(data, &rulesMap)
			if err != nil {
				return nil, err
			}
			err = toolbox.DefaultConverter.AssignConverted(&rules, rulesMap)
			return rules, err
		}
		rule := &Rule{}
		err := toolbox.DefaultConverter.AssignConverted(&rule, ruleMap)
		rules = append(rules, rule)
		return rules, err
	default:
		rule := &Rule{}
		if err := json.Unmarshal(data, rule); err != nil {
			err = json.Unmarshal(data, &rules)
			return rules, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
