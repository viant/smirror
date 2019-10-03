package smirror

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"smirror/config"
	"strings"
)

//Transfer represents a data transfer
type Transfer struct {
	Replace  map[string]string
	Resource *config.Resource
	Reader   io.Reader
	Dest     *Datafile
}

//GetReader returns a reader
func (t *Transfer) GetReader() (io.Reader, error) {
	if t.Reader == nil {
		return nil, fmt.Errorf("transfer reader was empty")
	}
	if len(t.Replace) > 0 {
		if err := t.replaceData(); err != nil {
			return nil, err
		}
	}
	return t.getReader()
}

func (t *Transfer) replaceData() error {
	data, err := ioutil.ReadAll(t.Reader)
	if err != nil {
		return err
	}
	textData := string(data)
	for k, v := range t.Replace {
		count := strings.Count(textData, k)
		if count == 0 {
			continue
		}
		textData = strings.Replace(textData, k, v, count)
	}
	t.Reader = strings.NewReader(textData)
	return nil
}

func (t *Transfer) getReader() (io.Reader, error) {
	if t.Dest.Compression == nil {
		return t.Reader, nil
	}
	if t.Dest.Codec == config.GZipCodec {
		buffer := new(bytes.Buffer)
		gzipWriter := gzip.NewWriter(buffer)
		_, err := io.Copy(gzipWriter, t.Reader)
		if err == nil {
			if err = gzipWriter.Flush(); err == nil {
				err = gzipWriter.Close()
			}
		}
		return buffer, err
	}
	return nil, fmt.Errorf("unsupported compression: %v", t.Dest.Compression.Codec)
}
