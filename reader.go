package mirror

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
)


//NewReader returns compression or regular reader
func NewReader(reader io.ReadCloser, compression *Compression) (io.ReadCloser, error) {
	if compression == nil {
		return reader, nil
	}
	payload, err:= ioutil.ReadAll(reader)
	_ = reader.Close()
	if err != nil {
		return nil, err
	}
	if compression.Codec == GZipCodec {
		return gzip.NewReader(bytes.NewReader(payload));
	}
	return nil, fmt.Errorf("unsupported code: %v", compression.Codec)
}