package mirror

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCompressionForURL(t *testing.T) {

	var useCases = []struct {
		description string
		URL         string
		expect      string
	}{
		{
			description: "gzip code",
			URL:         "ssh://asdsd/folder/file.txt.gz",
			expect:      GZipCodec,
		},
		{
			description: "empty code",
			URL:         "ssh://asdsd/folder/file.txt.x",
			expect:      "",
		},
	}

	for _, useCase := range useCases {
		actual := NewCompressionForURL(useCase.URL)
		if useCase.expect == "" {
			assert.Nil(t, actual, useCase.description)
			continue
		}
		if !assert.NotNil(t, actual, useCase.description) {
			continue
		}

		assert.Equal(t, useCase.expect, actual.Codec, useCase.description)

	}

}
