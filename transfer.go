package smirror

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"smirror/config"
	"smirror/transcoder"
)

//Transfer represents a data transfer
type Transfer struct {
	rule         *config.Rule
	partition    string
	splitCounter int32
	skipChecksum bool
	Resource     *config.Resource
	Reader       io.Reader
	Dest         *Datafile
}

//GetReader returns a reader
func (t *Transfer) GetReader() (reader io.Reader, err error) {
	if t.Reader == nil {
		return nil, fmt.Errorf("transfer reader was empty")
	}
	return t.getReader()
}

func (t *Transfer) getReader() (reader io.Reader, err error) {
	reader = t.Reader
	t.Reader = nil
	if t.Dest == nil {
		return reader, err
	}
	if t.rule.Transcoder != nil {
		reader, err = transcoder.NewReader(reader, t.rule.Transcoder, t.splitCounter)
		if err != nil {
			return nil, err
		}
	}

	if t.Dest.CompressionCodec() == config.GZipCodec {
		buffer := new(bytes.Buffer)
		gzipWriter := gzip.NewWriter(buffer)
		if _, err = io.Copy(gzipWriter, reader); err != nil {
			return nil, err
		}
		if err := gzipWriter.Flush(); err == nil {
			err = gzipWriter.Close()
		}
		return buffer, err
	}
	return reader, err
}
