package smirror

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"smirror/config"
	"strings"
	"unsafe"
)

const bufferSize = 1024 * 1024
var lineBreak = []byte{'\n'}

//Transfer represents a data transfer
type Transfer struct {
	Replace  []*config.Replace
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


func byteToString(data []byte) string {
	ptr := unsafe.Pointer(&data)
	return  *(*string)(ptr)
}



func (t *Transfer) replaceData() error {
	previous := ""
	replacer := t.Replacer()
	writer := new(bytes.Buffer)
	buffer := make([]byte, bufferSize)
	pending := make([]byte, bufferSize)
	for ; ; {
		bytesRead, e := t.Reader.Read(buffer)
		if bytesRead == 0 {
			if e != nil || e == io.EOF {
				break
			}
		}
		data := buffer[:bytesRead]
		for ;; {
			index := bytes.Index(data, lineBreak)
			if index == -1 {
				copy(pending, data)
				previous = byteToString(pending[:len(data)])
				break
			}
			text := byteToString(data[:index+1])
			if previous != "" {
				if _, err := replacer.WriteString(writer, previous + text); err != nil {
					return err
				}
				previous = ""
			} else if _, err := replacer.WriteString(writer, text); err != nil {
				return err
			}
			data = data[index+1:]
		}
	}
	t.Reader = writer
	return nil
}

func (t *Transfer) Replacer() (*strings.Replacer) {
	pairs := make([]string, 0)
	for _, replace := range t.Replace {
		pairs = append(pairs, replace.From)
		pairs = append(pairs, replace.To)
	}
	return strings.NewReplacer(pairs...)
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
