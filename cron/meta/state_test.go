package meta

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs/url"
	"testing"
	"time"
)

func TestState_Prune(t *testing.T) {

	var now = time.Now()

	var useCases = []struct {
		description string
		objects     map[string]time.Time
		maxAge      time.Duration
		expect      []string
	}{
		{
			description: "nothing to prune",
			objects: map[string]time.Time{
				"f1": now.Add(-1 * time.Second),
				"f2": now.Add(-10 * time.Second),
				"f3": now.Add(-20 * time.Second),
			},
			maxAge: 21 * time.Second,
			expect: []string{"f1", "f2", "f3"},
		},
		{
			description: "all pruned",
			objects: map[string]time.Time{
				"f1": now.Add(-5 * time.Second),
				"f2": now.Add(-10 * time.Second),
				"f3": now.Add(-20 * time.Second),
			},
			maxAge: 3 * time.Second,
			expect: []string{},
		},
		{
			description: "partial pruned",
			objects: map[string]time.Time{
				"f1": now.Add(-5 * time.Second),
				"f2": now.Add(-10 * time.Second),
				"f3": now.Add(-20 * time.Second),
			},
			maxAge: 10 * time.Second,
			expect: []string{"f1", "f2"},
		},
	}

	for i, useCase := range useCases {
		state := &State{}
		baseURL := fmt.Sprintf("mem://localhost/case%04d", i)
		state.Add(GetTestObjects(baseURL, useCase.objects)...)
		state.Prune(now, useCase.maxAge)
		resources := state.ProcessMap()
		assert.EqualValues(t, len(resources), len(useCase.expect), useCase.description)
		for _, name := range useCase.expect {
			URL := url.Join(baseURL, name)
			_, ok := resources[URL]
			assert.True(t, ok, useCase.description+" / "+URL)
		}
	}
}
