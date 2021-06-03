package endpoint

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/viant/afs"
	"os"
)

//Config represent  subscriber config
type Config struct {
	Queue       string
	BatchSize   int
	WaitTimeSec int64
}

//Initinitialises config
func (c *Config) Init(ctx context.Context, fs afs.Service) error {
	if c.BatchSize == 0 {
		c.BatchSize = 1
	}
	if c.WaitTimeSec == 0 {
		c.WaitTimeSec = 5
	}
	return nil
}

//Validate validates config
func (c *Config) Validate() error {
	if c.Queue == "" {
		return errors.New("Queue was empty")
	}
	return nil
}

func NewConfigFromEnv(key string) (*Config, error) {
	data := os.Getenv(key)
	cfg := &Config{}
	err := json.Unmarshal([]byte(data), cfg)
	return cfg, err
}

