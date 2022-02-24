package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/smirror/base"
	"testing"
)

func TestSplit_Name(t *testing.T) {

	var useCases = []struct {
		description string
		counter     int32
		URL         string
		split       *Split
		route       *Rule
		expect      string
	}{
		{
			description: "prefix template",
			route:       &Rule{},
			split: &Split{
				Template: "%03d_%s",
			},
			URL:    "gs://bucket/folder/data.csv.gz",
			expect: "folder/000_data.csv.gz",
		},
		{
			description: "suffix template",
			route:       &Rule{},
			counter:     2,
			split: &Split{
				Template: "%s_abc_%03d",
			},
			URL:    "gs://bucket/folder/data.csv.gz",
			expect: "folder/data_abc_002.csv.gz",
		},
		{
			description: "suffix template with 2 depth",
			route: &Rule{
				PreserveDepth: base.IntPtr(2),
			},
			counter: 32,
			split: &Split{
				Template: "%s_%03d",
			},
			URL:    "gs://bucket/folder1/subfolder/data.csv.gz",
			expect: "folder1/subfolder/data_032.csv.gz",
		},
	}

	for _, useCase := range useCases {
		actual := useCase.split.Name(useCase.route, useCase.URL, useCase.counter, nil)
		assert.EqualValues(t, useCase.expect, actual, useCase.description)
	}

}
