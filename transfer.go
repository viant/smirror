package smirror

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"smirror/config"
)

var lineBreak = []byte{'\n'}

//Transfer represents a data transfer
type Transfer struct {
	partition    string
	skipChecksum bool
	rewriter     *Rewriter
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

	if t.rewriter == nil {
		t.rewriter = NewRewriter()
	}

	if (t.Dest == nil || t.Dest.Compression == nil || t.Dest.Compression.Codec == "") && !t.rewriter.HasReplacer() {
		return reader, nil
	}

	buffer := new(bytes.Buffer)
	var writer io.Writer = buffer
	if t.Dest != nil && t.Dest.Compression != nil {
		if t.Dest.Compression.Codec != "" {
			if t.Dest.Codec == config.GZipCodec {
				gzipWriter := gzip.NewWriter(buffer)
				writer = gzipWriter
				defer func() {
					if err := gzipWriter.Flush(); err == nil {
						err = gzipWriter.Close()
					}
				}()
			}
		} else {
			return nil, fmt.Errorf("unsupported compression: %v", t.Dest.Compression.Codec)
		}
	}
	if !t.rewriter.HasReplacer() {
		_, err := io.Copy(writer, reader)
		return buffer, err
	}

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, bufferSize), 10*bufferSize)

	if scanner.Scan() {
		if err = t.rewriter.Write(writer, scanner.Bytes()); err != nil {
			return nil, err
		}
	}

	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return nil, err
		}
		if _, err = writer.Write(lineBreak); err != nil {
			return nil, err
		}
		if err = t.rewriter.Write(writer, scanner.Bytes()); err != nil {
			return nil, err
		}
	}
	return buffer, err
}
