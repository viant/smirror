package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs/matcher"
	"testing"
)

func TestRoutes_HasMatch(t *testing.T) {
	var useCases = []struct {
		description string
		Routes
		URL       string
		expectURL string
	}{
		{
			description: "suffix match",
			Routes: Routes{
				&Route{
					Basic: matcher.Basic{
						Suffix: ".tsv",
					},
					Dest: Resource{
						URL: "dst://abc",
					},
				},
				&Route{
					Basic: matcher.Basic{
						Suffix: ".csv",
					},
					Dest: Resource{
						URL: "dst://xyz",
					},
				},
			},

			URL:       "ssh://zz/folder/a.csv",
			expectURL: "dst://xyz",
		},
		{
			description: "prefix np match",
			Routes: Routes{
				&Route{
					Basic: matcher.Basic{
						Prefix: "/s",
					},
					Dest: Resource{
						URL: "dst://abc",
					},
				},
				&Route{
					Basic: matcher.Basic{
						Prefix: "/g",
					},
					Dest: Resource{
						URL: "dst://xyz",
					},
				},
			},

			URL:       "ssh://zz/folder/a.csv",
			expectURL: "",
		},
	}

	for _, useCase := range useCases {
		actual := useCase.HasMatch(useCase.URL)
		if useCase.expectURL == "" {
			assert.Nil(t, actual, useCase.description)
			continue
		}

		if !assert.NotNil(t, actual, useCase.description) {
			continue
		}

		assert.Equal(t, useCase.expectURL, actual.Dest.URL, useCase.description)
	}
}
