package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs/matcher"
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
				Basic: matcher.Basic{
					Prefix: "/folder/",
				},
			},
			URL:    "ssh:///folder/abc.xom",
			expect: true,
		},
		{
			description: "prefix no match",
			Route: Route{
				Basic: matcher.Basic{
					Prefix: "folder/",
				},
			},
			URL:    "ssh:///f/abc.xom",
			expect: false,
		},
		{
			description: "suffix match",
			Route: Route{
				Basic: matcher.Basic{
					Suffix: ".csv",
				},
			},
			URL:    "ssh:///folder/abc.csv",
			expect: true,
		},
		{
			description: "suffix no match",
			Route: Route{
				Basic: matcher.Basic{
					Suffix: ".tsv",
				},
			},
			URL:    "ssh:///f/abc.ts",
			expect: false,
		},
		{
			description: "filter no match",
			Route: Route{
				Basic: matcher.Basic{
					Suffix: ".tsv",
					Filter: `^[a-z]*/data/\\d+/`,
				},
			},
			URL:    "ssh://host/123/abc.tsv",
			expect: false,
		},
		{
			description: "filter match",
			Route: Route{
				Basic: matcher.Basic{
					Suffix: ".tsv",
					Filter: `^\/[a-z]+/data/\d+/`,
				},
			},
			URL:    "ssh://host/aa/data/002/abc.tsv",
			expect: true,
		},
	}

	for _, useCase := range useCases {
		actual := useCase.HasMatch(useCase.URL)
		assert.EqualValues(t, useCase.expect, actual, useCase.description)
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
