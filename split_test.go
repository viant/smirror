package mirror

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplit_Name(t *testing.T) {

	var useCases = []struct {
		description string
		counter     int32
		URL         string
		split       *Split
		route       *Route
		expect      string
	}{
		{
			description: "prefix template",
			route:       &Route{},
			split: &Split{
				Template: "%03d_%v",
			},
			URL:    "gs://bucket/folder/data.csv.gz",
			expect: "000_data.csv.gz",
		},
		{
			description: "suffix template",
			route:       &Route{},
			counter:     2,
			split: &Split{
				Template: "%v_abc_%03d",
			},
			URL:    "gs://bucket/folder/data.csv.gz",
			expect: "data_abc_002.csv.gz",
		},
		{
			description: "suffix template with 2 depth",
			route: &Route{
				FolderDepth: 2,
			},
			counter: 32,
			split: &Split{
				Template: "%v_%03d",
			},
			URL:    "gs://bucket/folder1/subfolder/data.csv.gz",
			expect: "folder1/subfolder/data_032.csv.gz",
		},
	}

	for _, useCase := range useCases {
		actual := useCase.split.Name(useCase.route, useCase.URL, useCase.counter)
		assert.EqualValues(t, useCase.expect, actual, useCase.description)
	}

}
