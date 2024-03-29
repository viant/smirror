package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs/matcher"
	"github.com/viant/smirror/base"
	"testing"
)

func TestRoute_HasMatch(t *testing.T) {

	var useCases = []struct {
		description string
		Rule
		URL    string
		expect bool
	}{
		{
			description: "prefix match",
			Rule: Rule{
				Source: &Resource{
					Basic: matcher.Basic{
						Prefix: "/folder/",
					},
				},
			},
			URL:    "ssh:///folder/abc.xom",
			expect: true,
		},
		{
			description: "prefix no match",
			Rule: Rule{
				Source: &Resource{
					Basic: matcher.Basic{
						Prefix: "folder/",
					},
				},
			},
			URL:    "ssh:///f/abc.xom",
			expect: false,
		},
		{
			description: "suffix match",
			Rule: Rule{
				Source: &Resource{
					Basic: matcher.Basic{
						Suffix: ".csv",
					},
				},
			},
			URL:    "ssh:///folder/abc.csv",
			expect: true,
		},
		{
			description: "suffix no match",
			Rule: Rule{
				Source: &Resource{
					Basic: matcher.Basic{
						Suffix: ".tsv",
					},
				},
			},
			URL:    "ssh:///f/abc.ts",
			expect: false,
		},
		{
			description: "filter no match",
			Rule: Rule{
				Source: &Resource{
					Basic: matcher.Basic{
						Suffix: ".tsv",
						Filter: `^[a-z]*/data/\\d+/`,
					},
				},
			},
			URL:    "ssh://host/123/abc.tsv",
			expect: false,
		},
		{
			description: "filter match",
			Rule: Rule{
				Source: &Resource{
					Basic: matcher.Basic{
						Suffix: ".tsv",
						Filter: `^\/[a-z]+/data/\d+/`,
					},
				},
			},
			URL:    "ssh://host/aa/data/002/abc.tsv",
			expect: true,
		},
		{
			description: "filter bucket match",
			Rule: Rule{
				Source: &Resource{
					Bucket: "xxx-billing",
					Basic: matcher.Basic{
						Suffix: ".zip",
						Prefix: `/yyyy-aws-billing-detailed-line-items-with-resources-and-tags-`,
					},
				},
			},
			URL:    "s3://xxx-billing/yyyy-aws-billing-detailed-line-items-with-resources-and-tags-2019-10.csv.zip",
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
		Rule
		URL    string
		expect string
	}{
		{
			description: "no folder depth",
			URL:         "s3://myducket/folder/asset1.txt",
			expect:      "folder/asset1.txt",
		},
		{
			description: "folder depth = 1",
			Rule: Rule{
				PreserveDepth: base.IntPtr(1),
			},
			URL:    "s3://myducket/folder/sub/asset1.txt",
			expect: "sub/asset1.txt",
		},
		{
			description: "folder depth = 2",
			Rule: Rule{
				PreserveDepth: base.IntPtr(2),
			},
			URL:    "s3://myducket/folder/sub/asset1.txt",
			expect: "folder/sub/asset1.txt",
		},
		{
			description: "folder depth exceeded path",
			Rule: Rule{
				PreserveDepth: base.IntPtr(4),
			},
			URL:    "s3://myducket/folder/sub/asset1.txt",
			expect: "folder/sub/asset1.txt",
		},
		{
			description: "gzip compression",
			Rule: Rule{
				PreserveDepth: base.IntPtr(0),
				Compression:   &Compression{Codec: GZipCodec},
			},
			URL:    "s3://myducket/folder/sub/asset1.txt",
			expect: "asset1.txt.gz",
		},
		{
			description: "folder root depth direction",
			Rule: Rule{
				PreserveDepth: base.IntPtr(-1),
			},
			URL:    "s3://myducket/folder/sub/dd/asset1.txt",
			expect: "sub/dd/asset1.txt",
		},
		{
			description: "folder root depth direction",
			Rule: Rule{
				PreserveDepth: base.IntPtr(-2),
			},
			URL:    "s3://myducket/folder/sub/dd/asset1.txt",
			expect: "dd/asset1.txt",
		},
		{
			description: "single asset from root",
			Rule: Rule{
				PreserveDepth: base.IntPtr(-1),
			},
			URL:    "s3://myducket/folder/asset1.txt",
			expect: "asset1.txt",
		},
		{
			description: "single asset from leaf",
			Rule: Rule{
				PreserveDepth: base.IntPtr(1),
			},
			URL:    "s3://myducket/folder/asset1.txt",
			expect: "folder/asset1.txt",
		},
	}

	for _, useCase := range useCases {
		actual := useCase.Name(useCase.URL)
		assert.EqualValues(t, useCase.expect, actual, useCase.description)
	}

}
