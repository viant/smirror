package smirror

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"smirror/config"
	"strings"
)

//Transfer represents a data transfer
type Transfer struct {
	Replace  []*config.Replace
	Size     int
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
	var lines = make([]string, 0)
	scanner := bufio.NewScanner(t.Reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, t.replace(scanner.Text()))
	}
	t.Reader = strings.NewReader(strings.Join(lines, "\n"))
	return nil
}

func (t *Transfer) replace(line string) string {
	for _, replace := range t.Replace {
		from := replace.From
		to := replace.To
		count := strings.Count(line, from)
		if count == 0 {
			continue
		}
		line = strings.Replace(line, from, to, count)
	}
	return line
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
