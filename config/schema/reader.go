package schema

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"smirror/base"
	"smirror/config"
	"smirror/contract"
	"strings"
	"unsafe"
)

const bufferSize = 1024 * 1024

var lineBreak = []byte{'\n'}

type reader struct {
	response   *contract.Response
	count      int32
	schema     *config.Schema
	replacer   *strings.Replacer
	scanner    *bufio.Scanner
	buf        *bytes.Buffer
	transient  *bytes.Buffer
	pending    int
	readEOF    bool
	badRecords int
	writeEOF   bool
}

func (t *reader) buffer() *bytes.Buffer {
	return t.buf
}

func (t *reader) transform() error {
	if !t.scanner.Scan() {
		t.readEOF = true
	}
	data := t.scanner.Bytes()
	var err error
	if t.replacer != nil {
		_, _ = t.replacer.WriteString(t.transient, byteToString(data))
		data, _ = ioutil.ReadAll(t.transient)
	}

	if len(data) == 0 {
		return nil
	}

	if t.schema != nil {
		if t.schema.IsJSON() && !json.Valid(data) {
			return t.reportBrokenJSON(data)
		}
		if t.schema.IsCSV() {
			data, err = t.adjustCSVValues(data)
		} else if t.schema.IsJSON() {
			data, err = t.adjustJSONValues(data, err)
		}
	}
	if err = t.failIfTooManyBadRecords(err); err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	if t.count > 0 {
		t.buf.Write(lineBreak)
		t.pending++
	}
	t.pending += len(data)
	t.buf.Write(data)
	t.count++
	return nil
}

func (t *reader) reportBrokenJSON(data []byte) (err error) {
	record := map[string]interface{}{}
	err = json.Unmarshal(data, &record)
	if err = t.failIfTooManyBadRecords(err); err != nil {
		return err
	}
	return nil
}

func (t *reader) adjustJSONValues(data []byte, err error) ([]byte, error) {
	if len(t.schema.Fields) == 0 {
		return data, nil
	}
	record := map[string]interface{}{}
	err = json.Unmarshal(data, &record)
	for _, field := range t.schema.Fields {
		value, ok := record[field.Name]
		if ! ok {
			continue
		}
		if record[field.Name], err = field.AdjustValue(value); err != nil {
			return nil, err
		}
	}
	if updated, err := json.Marshal(record); err == nil {
		data = updated
	}
	return data, err
}

func (t *reader) failIfTooManyBadRecords(err error) error {
	if err == nil || t.schema.MaxBadRecords == nil {
		return nil
	}
	t.badRecords++
	if t.response != nil {
		t.response.BadRecords++
	}
	if t.badRecords >= *t.schema.MaxBadRecords {
		return base.NewSchemaError(errors.Wrapf(err, "too many bad records: %v", t.badRecords))
	}
	return nil
}

func (t *reader) adjustCSVValues(data []byte) ([]byte, error) {
	csvReader := t.schema.NewCsvReader(bytes.NewReader(data))
	record, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	if t.schema.FieldCount > 0 {
		if len(record) > t.schema.FieldCount {
			record = record[:t.schema.FieldCount]
		} else {
			for len(record) < t.schema.FieldCount {
				record = append(record, "")
			}
		}
	}
	if err = t.adjustCSVDataType(record); err != nil {
		return nil, err
	}

	writer := csv.NewWriter(t.transient)
	writer.Comma = csvReader.Comma
	writer.UseCRLF = false
	if err = writer.Write(record); err != nil {
		return nil, err
	}
	writer.Flush()
	data, _ = ioutil.ReadAll(t.transient)
	if len(data) == 0 {
		return data, nil
	}
	return data[:len(data)-1], err
}

func (t *reader) Read(p []byte) (n int, err error) {
	if t.writeEOF {
		return 0, io.EOF
	}
	expect := len(p)
	for t.pending < expect && !t.readEOF {
		err := t.transform()
		if err != nil {
			return 0, err
		}
	}

	read, err := t.buffer().Read(p)
	if err == io.EOF || read == 0 {
		if t.readEOF {
			t.writeEOF = true
		} else {
			err = nil
		}
	}
	t.pending -= read
	return read, err
}

func (t *reader) adjustCSVDataType(record []string) error {
	if len(t.schema.Fields) == 0 {
		return nil
	}
	var err error
	for _, field := range t.schema.Fields {
		if field.Position == nil {
			continue
		}
		index := *field.Position
		if index >= len(record) {
			continue
		}
		record[index], err = field.AdjustText(record[index])
		if err != nil {
			return err
		}
	}
	return err
}

func NewReader(r io.Reader, rule *config.Rule, response *contract.Response) (io.Reader, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, bufferSize), 10*bufferSize)
	return &reader{
		response:  response,
		schema:    rule.Schema,
		transient: new(bytes.Buffer),
		buf:       new(bytes.Buffer),
		replacer:  rule.NewReplacer(),
		scanner:   scanner,
	}, nil

}

func byteToString(data []byte) string {
	ptr := unsafe.Pointer(&data)
	return *(*string)(ptr)
}
