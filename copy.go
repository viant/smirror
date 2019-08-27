package smirror

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"smirror/config"
)

//Copy represents a data copy
type Copy struct {
	Resource *config.Resource
	Reader io.Reader
	Dest   *Datafile
}

//GetReader returns a reader
func (c *Copy) GetReader() (io.Reader, error) {
	if c.Reader == nil {
		return nil, fmt.Errorf("copy reader was empty")
	}
	if c.Dest.Compression == nil {
		return c.Reader, nil
	}
	if c.Dest.Codec == config.GZipCodec {

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
