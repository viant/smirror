package smirror

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"smirror/config"
	"testing"
)

func TestNewReader(t *testing.T) {

	var useCases = []struct {
		description string
		data        string
		compression *config.Compression
		hasError    bool
	}{

		{
			description: "raw data",
			data:        "abc",
		},

		{
			description: "compressed data",
			data:        "abcd",
			compression: &config.Compression{Codec: config.GZipCodec},
		},
		{
			description: "unknown code error",
			data:        "abcd",
			compression: &config.Compression{Codec: "abc"},
			hasError:    true,
		},
	}

	for _, useCase := range useCases {

		data := []byte(useCase.data)
		if useCase.compression != nil {
			switch useCase.compression.Codec {
			case config.GZipCodec:
				buffer := new(bytes.Buffer)
				writer := gzip.NewWriter(buffer)
				_, _ = writer.Write(data)
				_ = writer.Flush()
				_ = writer.Close()
				data = buffer.Bytes()
			}
		}

		reader, err := NewReader(ioutil.NopCloser(bytes.NewReader(data)), useCase.compression)
		if useCase.hasError {
			assert.NotNil(t, err, useCase.description)
			continue
		}
		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		data, err = ioutil.ReadAll(reader)
		assert.Nil(t, err, useCase.description)
		assert.Equal(t, useCase.data, string(data), useCase.description)
		assert.Nil(t, reader.Close(), useCase.description)
	}

}
