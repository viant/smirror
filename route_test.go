package smirror

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute_HasMatch(t *testing.T) {

	var useCases = []struct {
		description string
		Route
		URL    string
		expect bool
	}{
		{
			description: "prefix match",
			Route: Route{
				Prefix: "/folder/",
			},
			URL:    "ssh:///folder/abc.xom",
			expect: true,
		},
		{
			description: "prefix no match",
			Route: Route{
				Prefix: "folder/",
			},
			URL:    "ssh:///f/abc.xom",
			expect: false,
		},
		{
			description: "suffix match",
			Route: Route{
				Suffix: ".csv",
			},
			URL:    "ssh:///folder/abc.csv",
			expect: true,
		},
		{
			description: "suffix no match",
			Route: Route{
				Suffix: ".tsv",
			},
			URL:    "ssh:///f/abc.ts",
			expect: false,
		},
		{
			description: "filter no match",
			Route: Route{
				Suffix: ".tsv",
				Filter:`^[a-z]*/data/\\d+/`,
			},
			URL:    "ssh://host/123/abc.tsv",
			expect: false,
		},
		{
			description: "filter match",
			Route: Route{
				Suffix: ".tsv",
				Filter:`^\/[a-z]+/data/\d+/`,
			},
			URL:    "ssh://host/aa/data/002/abc.tsv",
			expect: true,
		},

	}

	for _, useCase := range useCases {
		err := useCase.Init()
		assert.Nil(t, err, useCase.description)
		actual := useCase.HasMatch(useCase.URL)
		assert.EqualValues(t, useCase.expect, actual, useCase.description)
	}

}

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
					Suffix:  ".tsv",
					DestURL: "dst://abc",
				},
				&Route{
					Suffix:  ".csv",
					DestURL: "dst://xyz",
				},
			},

			URL:       "ssh://zz/folder/a.csv",
			expectURL: "dst://xyz",
		},
		{
			description: "prefix np match",
			Routes: Routes{
				&Route{
					Prefix:  "/s",
					DestURL: "dst://abc",
				},
				&Route{
					Prefix:  "/g",
					DestURL: "dst://xyz",
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

		assert.Equal(t, useCase.expectURL, actual.DestURL, useCase.description)
	}
}

func TestRoute_Name(t *testing.T) {

	var useCases = []struct {
		description string
		Route
		URL    string
		expect string
	}{
		{
			description: "no folder depth",
			URL:         "s3://myducket/folder/asset1.txt",
			expect:      "asset1.txt",
		},
		{
			description: "folder depth = 1",
			Route: Route{
				FolderDepth: 1,
			},
			URL:    "s3://myducket/folder/sub/asset1.txt",
			expect: "sub/asset1.txt",
		},
		{
			description: "folder depth = 2",
			Route: Route{
				FolderDepth: 2,
			},
			URL:    "s3://myducket/folder/sub/asset1.txt",
			expect: "folder/sub/asset1.txt",
		},
		{
			description: "folder depth exceeded path",
			Route: Route{
				FolderDepth: 4,
			},
			URL:    "s3://myducket/folder/sub/asset1.txt",
			expect: "folder/sub/asset1.txt",
		},
		{
			description: "gzip compression",
			Route: Route{
				Compression: &Compression{Codec: GZipCodec},
			},
			URL:    "s3://myducket/folder/sub/asset1.txt",
			expect: "asset1.txt.gz",
		},
	}

	for _, useCase := range useCases {
		actual := useCase.Name(useCase.URL)
		assert.EqualValues(t, useCase.expect, actual, useCase.description)
	}

}
