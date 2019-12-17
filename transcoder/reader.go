package transcoder

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/actgardner/gogen-avro/container"
	"github.com/pkg/errors"
	"github.com/viant/toolbox/data"
	"io"
	"smirror/config"
	"smirror/transcoder/avro"
	"smirror/transcoder/avro/schma"
)

const (
	bufferSize = 1024 * 1024
)

type reader struct {
	*config.Transcoding
	reader     io.Reader
	scanner    *bufio.Scanner
	buffer     *bytes.Buffer
	avroRecord *avro.Record
	avroWriter *container.Writer
	record     map[string]interface{}
	eof        bool
	pending    bool
}

func (t *reader) Read(p []byte) (n int, err error) {
	if t.eof {
		return 0, io.EOF
	}

	for t.pending && !t.eof {
		err := t.transform()
		if err != nil {
			return 0, err
		}
	}
	read, err := t.buffer.Read(p)
	if err == io.EOF && !t.eof {
		t.pending = true
		err = nil
	}

	return read, err
}

func (t *reader) next() error {
	data := t.scanner.Bytes()
	if t.Source.IsCSV() {
		reader := t.csvReader(bytes.NewReader(data))
		values, err := reader.Read()
		if err != nil {
			return err
		}
		for i := range t.Source.Fields {
			if i >= len(values) {
				break
			}
			t.record[t.Source.Fields[i]] = values[i]
		}
		return nil
	}
	if !t.Source.IsJSON() {
		return errors.Errorf("unsupported source format: %v", t.Source.Format)
	}
	return json.Unmarshal(data, &t.record)
}

func (t *reader) nextRecord() error {
	err := t.next()
	if err != nil || len(t.Transcoding.PathMapping) == 0 {
		return err
	}

	original := data.Map(t.record)
	mapped := data.NewMap()
	for _, mapping := range t.Transcoding.PathMapping {
		value, ok := original.GetValue(mapping.From)
		if !ok {
			continue
		}
		mapped.SetValue(mapping.To, value)
	}
	t.record = mapped
	return err
}

func (t *reader) transform() error {
	if t.Dest.IsAvro() {
		return t.transformToAVRO()
	}
	return fmt.Errorf("unsupported avro format")
}

func (t *reader) transformToAVRO() error {

	for !t.eof {
		err := t.transformRecordToAVRO()
		if err != nil {
			return err
		}
	}
	t.pending = t.buffer.Len() > 0
	return nil
}

func (t *reader) transformRecordToAVRO() error {
	hasMore := t.scanner.Scan()
	if !hasMore {
		t.eof = true
		return t.avroWriter.Flush()
	}
	if err := t.nextRecord(); err != nil {
		return err
	}
	t.avroRecord.Data = t.record
	return t.avroWriter.WriteRecord(t.avroRecord)
}

func (t *reader) csvReader(reader io.Reader) *csv.Reader {
	result := csv.NewReader(reader)
	if t.Source.Delimiter != "" {
		result.Comma = rune(t.Source.Delimiter[0])
	}
	result.LazyQuotes = t.Source.LazyQuotes
	return result
}

//NewReader creates a transcoding reader
func NewReader(r io.Reader, transcoding *config.Transcoding) (io.Reader, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, bufferSize), 10*bufferSize)
	result := &reader{
		Transcoding: transcoding,
		reader:      r,
		scanner:     scanner,
		record:      make(map[string]interface{}),
		pending:     true,
	}
	if transcoding.Dest.IsAvro() {
		result.buffer = new(bytes.Buffer)
		rawSchema := transcoding.Dest.Schema
		if rawSchema == "" {
			return nil, errors.Errorf("avro schema was empty, %v", transcoding.Dest.SchemaURL)
		}
		avroSchema, err := schma.New(rawSchema)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to load schema: %v", transcoding.Dest.SchemaURL)
		}
		err = avro.SetWriter(avroSchema)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to set avro writer: %v", transcoding.Dest.SchemaURL)
		}
		result.avroRecord = avro.NewRecord(result.record, avroSchema, rawSchema)
		result.avroWriter, err = container.NewWriter(result.buffer, container.Snappy, 40, rawSchema)
	}
	return result, nil
}
