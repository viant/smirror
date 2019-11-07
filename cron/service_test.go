package cron

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs"
	"github.com/viant/afs/asset"
	"github.com/viant/afs/file"
	"github.com/viant/afs/mem"
	"github.com/viant/afs/url"
	"github.com/viant/assertly"
	"os"
	"smirror/base"
	"smirror/cron/meta"
	"strings"
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
          "Source": {
		  	  "URL": "mem://localhost/case001/"
          },
		  "Dest": {
  			  "URL": "mem://localhost/"
          }
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

		if !assert.Equal(t, base.StatusOK, response.Status, useCase.description+" "+response.Error) {
			continue
		}
		if !assert.Equal(t, 1, len(response.Matched), useCase.description) {
			continue
		}
		actual, err := loadMeta(ctx, fs, useCase.config.MetaURL)
		assertly.AssertValues(t, useCase.expect, actual, useCase.description)

		for i := range resources {

			sourceURL := url.Join(useCase.baseURL, resources[i].Name)
			basePath := strings.Replace(sourceURL, "mem://", "", 1)
			destURL := url.Join(useCase.config.Resources.Rules[0].Dest.URL, basePath)
			exists, _ := fs.Exists(ctx, destURL)
			assert.True(t, exists)
			_ = fs.Delete(ctx, destURL)

		}
		response = service.Tick(ctx)
		if !assert.Equal(t, 0, len(response.Matched), useCase.description) {
			continue
		}
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
