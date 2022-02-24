package transcoder

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/actgardner/gogen-avro/v10/container"
	"github.com/pkg/errors"
	"github.com/viant/smirror/config"
	"github.com/viant/smirror/transcoder/avro"
	"github.com/viant/smirror/transcoder/avro/schma"
	"github.com/viant/smirror/transcoder/xlsx"
	"github.com/viant/toolbox/data"
	"io"
)

const (
	bufferSize = 1024 * 1024
)

type reader struct {
	badRecords int
	*config.Transcoding
	splitCounter int32
	fileds       []string
	reader       io.Reader
	scanner      *bufio.Scanner
	xlsDecoder   *xlsx.Decoder
	output       *bytes.Buffer
	avroRecord   *avro.Record
	avroWriter   *container.Writer
	record       map[string]interface{}
	count        int
	readEOF      bool
	writeEOF     bool
	pending      bool
}

func (t *reader) Read(p []byte) (int, error) {
	if t.writeEOF {
		return 0, io.EOF
	}
	for t.pending && !t.readEOF {
		err := t.transform()
		if err != nil {
			return 0, err
		}
	}
	read, err := t.output.Read(p)
	if err == io.EOF || read == 0 {
		if t.readEOF {
			t.writeEOF = true
		} else {
			t.pending = true
			err = nil
		}
	}
	return read, err
}

func (t *reader) next() error {

	if t.xlsDecoder != nil {
		return t.xlsDecoder.NextRecord(t.record)
	}
	line := t.scanner.Bytes()
	if t.Source.HasHeader && t.splitCounter == 0 {
		reader := t.csvReader(bytes.NewReader(line))
		fields, err := reader.Read()
		if err != nil {
			return err
		}
		if len(t.Source.Fields) == 0 {
			t.Source.Fields = fields
			t.fileds = fields
		}
		if !t.scanner.Scan() {
			return nil
		}
		line = t.scanner.Bytes()
	}

	if t.Source.IsCSV() {
		reader := t.csvReader(bytes.NewReader(line))
		values, err := reader.Read()
		if err != nil {
			return err
		}
		for i := range t.fileds {
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
	return json.Unmarshal(line, &t.record)
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
	if t.Dest.IsJSON() {
		return t.transformToJSON()
	}
	if t.Dest.IsAvro() {
		return t.transformToAVRO()
	}
	return fmt.Errorf("unsupported %v dest format", t.Dest.Format)
}

func (t *reader) transformToAVRO() error {
	var err error
	for !t.readEOF {
		err = t.transformRecordToAVRO()
		if err != nil || t.count%int(t.Dest.RecordPerBlock) == 0 {
			return err
		}
	}
	t.pending = t.output.Len() > 0
	return err
}

func (t *reader) hasMore() bool {
	if t.scanner != nil {
		return t.scanner.Scan()
	}
	return t.xlsDecoder.HasMore()
}

func (t *reader) transformRecordToAVRO() error {
	hasMore := t.hasMore()
	if !hasMore {
		t.readEOF = true
		return t.avroWriter.Flush()
	}
	if err := t.nextRecord(); err != nil {
		t.badRecords++
		if t.Transcoding.MaxBadRecords != nil && t.badRecords >= *t.Transcoding.MaxBadRecords {
			return errors.Wrapf(err, "too many bad records: %v, max allowed: %v", t.badRecords, *t.Transcoding.MaxBadRecords)
		}
	}
	if len(t.record) == 0 {
		return nil
	}
	t.avroRecord.Data = t.record
	t.count++
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

func (t *reader) transformToJSON() error {
	var err error
	for !t.readEOF {
		err = t.transformRecordToJSON()
		if err != nil {
			return err
		}
	}
	t.pending = t.output.Len() > 0
	return err

}

func (t *reader) transformRecordToJSON() error {
	hasMore := t.hasMore()
	if !hasMore {
		t.readEOF = true
		return nil
	}

	if err := t.nextRecord(); err != nil {
		t.badRecords++
		if t.Transcoding.MaxBadRecords != nil && t.badRecords >= *t.Transcoding.MaxBadRecords {
			return errors.Wrapf(err, "too many bad records: %v, max allowed: %v", t.badRecords, *t.Transcoding.MaxBadRecords)
		}
	}
	if len(t.record) == 0 {
		return nil
	}

	data, err := json.Marshal(t.record)
	if err != nil {
		t.badRecords++
		if t.Transcoding.MaxBadRecords != nil && t.badRecords >= *t.Transcoding.MaxBadRecords {
			return errors.Wrapf(err, "too many bad records: %v, max allowed: %v", t.badRecords, *t.Transcoding.MaxBadRecords)
		}
	}
	if t.output.Len() > 0 {
		t.output.WriteString("\n")
	}
	_, err =  t.output.Write(data)
	return err
}

//NewReader creates a transcoding reader
func NewReader(r io.Reader, transcoding *config.Transcoding, splitCounter int32) (io.Reader, error) {
	var err error
	var scanner *bufio.Scanner
	var xlsDecoder *xlsx.Decoder
	if transcoding.Source.IsCSV() || transcoding.Source.IsJSON() {
		scanner = bufio.NewScanner(r)
		scanner.Buffer(make([]byte, bufferSize), 10*bufferSize)
	} else if transcoding.Source.IsXLSX() {
		xlsDecoder, err = xlsx.NewDecoder(r)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open xls decoder")
		}
	}
	result := &reader{
		xlsDecoder:   xlsDecoder,
		splitCounter: splitCounter,
		Transcoding:  transcoding,
		reader:       r,
		scanner:      scanner,
		record:       make(map[string]interface{}),
		fileds:       transcoding.Source.Fields,
		pending:      true,
	}

	if transcoding.Dest.IsJSON() {
		result.output = new(bytes.Buffer)
		return result, nil
	}

	if transcoding.Dest.IsAvro() {
		result.output = new(bytes.Buffer)
		rawSchema := transcoding.Dest.Schema

		if transcoding.Source.IsXLSX() {
			rawSchema = xlsDecoder.Schema()
		}
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

		if transcoding.Dest.RecordPerBlock == 0 {
			transcoding.Dest.RecordPerBlock = 20
		}
		recordPerBlock := transcoding.Dest.RecordPerBlock
		result.avroWriter, err = container.NewWriter(result.output, container.Snappy, recordPerBlock, rawSchema)

	}
	return result, nil
}
