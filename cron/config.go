package cron

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/viant/smirror/base"
	"github.com/viant/smirror/cron/config"
	"github.com/viant/afs"
	"github.com/viant/toolbox"
	"os"
	"strings"
)

//Config represents cron config
type Config struct {
	base.Config
	MetaURL    string
	TimeWindow config.TimeWindow
	Resources  config.Ruleset
}

//Load initialises routes
func (c *Config) Init(ctx context.Context, fs afs.Service) error {
	c.Config.Init()
	c.TimeWindow.Init()
	if err := c.TimeWindow.Validate(); err != nil {
		return err
	}
	if c.MetaURL == "" {
		return errors.New("metaURL was empty")
	}
	return c.Resources.Init(ctx, fs, c.ProjectID)
}

//NewConfigFromEnv returns new config from env
func NewConfigFromEnv(ctx context.Context, key string) (*Config, error) {
	JSONOrURL := strings.TrimSpace(os.Getenv(key))
	if toolbox.IsStructuredJSON(JSONOrURL) {
		return NewConfigFromJSON(ctx, JSONOrURL)
	}
	return NewConfigFromURL(ctx, JSONOrURL)
}

//NewConfigFromJSON creates a new config from env
func NewConfigFromJSON(ctx context.Context, payload string) (*Config, error) {
	cfg := &Config{}
	err := json.NewDecoder(strings.NewReader(payload)).Decode(cfg)
	if err == nil {
		err = cfg.Init(ctx, afs.New())
	}

	return cfg, err
}

//NewConfigFromURL creates a new config from env
func NewConfigFromURL(ctx context.Context, URL string) (*Config, error) {
	service := afs.New()
	reader, err := service.OpenURL(ctx, URL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download config: %v", URL)
	}
	cfg := &Config{}
	err = json.NewDecoder(reader).Decode(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode config: %v ", URL)
	}
	return cfg, cfg.Init(ctx, afs.New())
}
