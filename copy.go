package smirror

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

//Copy represents a data copy
type Copy struct {
	Reader io.Reader
	Dest   *Datafile
}

//GetReader returns a reader
func (c *Copy) GetReader() (io.Reader, error) {
	if c.Dest.Compression == nil {
		return c.Reader, nil
	}
	if c.Dest.Codec == GZipCodec {

		buffer := new(bytes.Buffer)
		gzipWriter := gzip.NewWriter(buffer)
		_, err := io.Copy(gzipWriter, c.Reader)
		if err == nil {
			if err = gzipWriter.Flush(); err == nil {
				err = gzipWriter.Close()
			}
		}
		return buffer, err
	}
	return nil, fmt.Errorf("unsupported compression: %v", c.Dest.Compression.Codec)
}
