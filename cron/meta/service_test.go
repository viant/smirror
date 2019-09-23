package meta

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs"
	"github.com/viant/afs/url"
	"testing"
	"time"
)

func TestService_PendingResources(t *testing.T) {


	now := time.Now()
	var useCases = []struct {
		description    string
		baseURL        string
		processedSoFar map[string]time.Time
		candidates     map[string]time.Time
		expect         []string
	}{
		{
			description:"no pending resources",
			processedSoFar:map[string]time.Time{
				"f1":now.Add(-15 *time.Second),
				"f2":now.Add(-12 *time.Second),
				"f3":now.Add(-10 *time.Second),
			},
			candidates:map[string]time.Time{
				"f2":now.Add(-12 *time.Second),
				"f3":now.Add(-10 *time.Second),
			},
		},
		{
			description:"all pending resources",
			processedSoFar:map[string]time.Time{
				"f1":now.Add(-15 *time.Second),
			},
			candidates:map[string]time.Time{
				"f2":now.Add(-12 *time.Second),
				"f3":now.Add(-10 *time.Second),
			},
			expect:[]string{"f2", "f3"},
		},
		{
			description:"partial pending resources",
			processedSoFar:map[string]time.Time{
				"f1":now.Add(-15 *time.Second),
			},
			candidates:map[string]time.Time{
				"f1":now.Add(-15 *time.Second),
				"f2":now.Add(-12 *time.Second),
				"f3":now.Add(-10 *time.Second),
			},
			expect:[]string{"f2", "f3"},
		},

	}

	ctx := context.Background()
	fs := afs.New()
	for i, useCase := range useCases {
		if useCase.baseURL == "" {
			useCase.baseURL = fmt.Sprintf("mem://localhost/case%04d", i)
		}
		service := New(url.Join(useCase.baseURL, "/meta.json"), 0, fs)
		processed := GetTestObjects(useCase.baseURL, useCase.processedSoFar)
		err := service.AddProcessed(ctx, processed)
		assert.Nil(t, err, useCase.description)
		candidates := GetTestObjects(useCase.baseURL, useCase.candidates)
		pending, err := service.PendingResources(ctx, candidates)
		assert.Nil(t, err, useCase.description)
		assert.EqualValues(t, len(pending), len(useCase.expect))
		actual := map[string]bool{}
		for _, elem := range pending {
			actual[elem.URL()] = true
		}
		for _, expect := range useCase.expect {
			URL := url.Join(useCase.baseURL, expect)
			_, ok := actual[URL]
			assert.True(t, ok, useCase.description+" / "+URL)
		}
	}
}
