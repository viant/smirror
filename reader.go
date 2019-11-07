package smirror

import (
	"compress/gzip"
	"fmt"
	"io"
	"smirror/config"
)

//NewReader returns compression or regular reader returning raw data
func NewReader(reader io.ReadCloser, compression *config.Compression) (io.ReadCloser, error) {
	if compression == nil || !compression.Uncompress {
		return reader, nil
	}
	if compression.Codec == config.GZipCodec {
		return gzip.NewReader(reader)
	}
	return nil, fmt.Errorf("unsupported code: %v", compression.Codec)
}
