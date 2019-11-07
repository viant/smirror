package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs/matcher"
	"testing"
)

func TestRoutes_HasMatch(t *testing.T) {
	var useCases = []struct {
		description string
		Ruleset
		URL       string
		expectURL string
	}{
		{
			description: "suffix match",
			Ruleset: Ruleset{
				Rules: []*Rule{
					{
						Source: &Resource{
							Basic: matcher.Basic{
								Suffix: ".tsv",
							},
						},
						Dest: &Resource{
							URL: "dst://abc",
						},
					},
					{
						Source: &Resource{
							Basic: matcher.Basic{
								Suffix: ".csv",
							},
						},
						Dest: &Resource{
							URL: "dst://xyz",
						},
					},
				},
			},
			URL:       "ssh://zz/folder/a.csv",
			expectURL: "dst://xyz",
		},
		{
			description: "prefix np match",
			Ruleset: Ruleset{
				Rules: []*Rule{
					{
						Source: &Resource{
							Basic: matcher.Basic{
								Prefix: "/s",
							},
						},
						Dest: &Resource{
							URL: "dst://abc",
						},
					},
					{
						Source: &Resource{
							Basic: matcher.Basic{
								Prefix: "/g",
							},
						},
						Dest: &Resource{
							URL: "dst://xyz",
						},
					},
				},
			},

			URL:       "ssh://zz/folder/a.csv",
			expectURL: "",
		},
	}

	for _, useCase := range useCases {
		matched := useCase.Match(useCase.URL)
		var actual *Rule
		if len(matched) == 1 {
			actual = matched[0]
		}

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
