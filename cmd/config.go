package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	"smirror"
	"smirror/base"
	"strings"
)

var async = false

var configURL = url.Join(base.InMemoryStorageBaseURL, "/Smirror/config/")
var ruleBaseURL = url.Join(configURL, "rule")

//NewConfig creates bqtail config
func NewConfig(ctx context.Context, projectID string) (*smirror.Config, error) {
	cfg, err := newConfig(ctx, projectID)
	if err != nil {
		return cfg, err
	}
	configJSON, _ := json.Marshal(cfg)
	fs := afs.New()
	if err := fs.Upload(ctx, cfg.URL, file.DefaultFileOsMode, bytes.NewReader(configJSON)); err != nil {
		return nil, errors.Wrapf(err, "failed to upload config: %v", cfg.URL)
	}
	emptyRuleURL := url.Join(ruleBaseURL, "t")
	_ = fs.Upload(ctx, emptyRuleURL, file.DefaultFileOsMode, strings.NewReader("."))
	err = cfg.Init(ctx, fs)
	return cfg, err
}

func newConfig(ctx context.Context, projectID string) (*smirror.Config, error) {
	var err error
	cfg := &smirror.Config{}
	cfg.ProjectID = projectID
	cfg.Mirrors.BaseURL = ruleBaseURL
	cfg.Mirrors.CheckInMs = 0
	cfg.MaxRetries = 3
	cfg.URL = url.Join(configURL, "config.json")
	return cfg, err
}
