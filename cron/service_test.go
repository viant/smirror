package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs"
	"github.com/viant/afs/asset"
	"github.com/viant/afs/file"
	"github.com/viant/afs/mem"
	"github.com/viant/afs/url"
	"github.com/viant/assertly"
	"os"
	"smirror/cron/meta"
	"testing"
	"time"
)

func TestService_Tick(t *testing.T) {

	now := time.Now()
	fs := afs.New()
	var useCases = []struct {
		description   string
		configKey     string
		configPayload string
		config        *Config
		baseURL       string
		input         map[string]time.Time
		expect        interface{}
	}{
		{
			description: "trigger 1 events",
			configPayload: `{
  "MetaURL": "mem://localhost/ops/meta.json",
  "TimeWindow": {
    "DurationInSec": 5
  },
  "Resources": {	
	  "Rules": [
		{
		  "URL": "mem://localhost/case001/",
		  "DestFunction": "func1"
		}
	  ]
  }
}
`,
			configKey: "CONFIG",
			baseURL:   "mem://localhost/case001/",
			input: map[string]time.Time{
				"f1.txt": now.Add(-1 * time.Millisecond),
			},
			expect: `{"Processed":[{"URL":"mem://localhost/case001/f1.txt"}]}`,
		},
	}

	ctx := context.Background()
	for _, useCase := range useCases {
		resources := getTestObjects(useCase.input)
		err := asset.Create(mem.Singleton(), useCase.baseURL, resources)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		_ = os.Setenv(useCase.configKey, useCase.configPayload)
		useCase.config, err = NewConfigFromEnv(ctx, useCase.configKey)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		service, err := New(ctx, useCase.config, fs)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		response := service.Tick(ctx)
		if !assert.Nil(t, response.Error != "", useCase.description) {
			continue
		}

		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		actual, err := loadMeta(ctx, fs, useCase.config.MetaURL)
		assertly.AssertValues(t, useCase.expect, actual, useCase.description)

		funcURL := fmt.Sprintf("mem://localhost/%v", useCase.config.Resources.Rules[0].DestFunction)

		exists, _ := fs.Exists(ctx, funcURL)
		assert.True(t, exists)
		_ = fs.Delete(ctx, funcURL)

		response = service.Tick(ctx)
		if !assert.Nil(t, response.Error != "", useCase.description) {
			continue
		}

		exists, _ = fs.Exists(ctx, funcURL)
		//make sure that function is triggered only once
		assert.False(t, exists, useCase.description)

	}

}

func loadMeta(ctx context.Context, fs afs.Service, URL string) (*meta.State, error) {
	reader, err := fs.DownloadWithURL(ctx, URL)
	if err != nil {
		return nil, err
	}
	result := &meta.State{}
	return result, json.NewDecoder(reader).Decode(result)
}

func getTestObjects(objects map[string]time.Time) []*asset.Resource {
	var result = make([]*asset.Resource, 0)
	for name, mod := range objects {
		_, name := url.Split(name, mem.Scheme)
		info := file.NewInfo(name, 0, 0644, mod, false)
		result = append(result, asset.NewFile(name, []byte("test"), info.Mode()))
	}
	return result
}
