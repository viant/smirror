package endpoint

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/viant/afs"
	"os"
)

//Config represent  client config
type Config struct {
	Subscription string
	ProjectID    string
	BatchSize    int
	WaitTimeSec  int64
	VisibilityTimeout int64
}

//Initinitialises config
func (c *Config) Init(ctx context.Context, fs afs.Service) error {
	if c.BatchSize == 0 {
		c.BatchSize = 1
	}
	if c.WaitTimeSec == 0 {
		c.WaitTimeSec = 5
	}
	if c.VisibilityTimeout == 0 {
		c.VisibilityTimeout = 60
	}
	return nil
}

//Validate validates config
func (c *Config) Validate() error {
	if c.Subscription == "" {
		return errors.New("subscription was empty")
	}
	return nil
}

func NewConfigFromEnv(key string) (*Config, error) {
	data := os.Getenv(key)
	cfg := &Config{}
	err := json.Unmarshal([]byte(data), cfg)
	return cfg, err
}
