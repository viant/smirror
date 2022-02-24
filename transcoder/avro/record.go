package avro

import (
	"github.com/viant/toolbox"
	"io"
	"github.com/viant/smirror/transcoder/avro/schma"
)

func translateToRecord(schema *schma.Schema) schma.Translator {
	return func(value interface{}, w io.Writer) error {
		record := toolbox.AsMap(value)
		for _, field := range schema.Fields {
			val := record[field.Name]
			err := field.Type.Write(val, w)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

//Record represents an avro map based record
type Record struct {
	schema     *schma.Schema
	schemaText string
	Data       map[string]interface{}
}

func (r Record) Serialize(w io.Writer) error {
	return r.schema.Write(r.Data, w)
}

func (r Record) Schema() string {
	return r.schemaText
}

//NewRecord create a avro record
func NewRecord(record map[string]interface{}, schema *schma.Schema, schemaText string) *Record {
	return &Record{
		schema:     schema,
		schemaText: schemaText,
		Data:       record,
	}
}
