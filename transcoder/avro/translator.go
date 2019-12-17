package avro

import (
	"github.com/viant/toolbox"
	"smirror/transcoder/avro/schma"
)

//TranslateRecord translates record according avro schema
func TranslateRecord(input map[string]interface{}, avroSchema *schma.Schema) map[string]interface{} {
	result := translate(input, avroSchema)
	return toolbox.AsMap(result)
}

func translate(input map[string]interface{}, avroSchema *schma.Schema) interface{} {
	if avroSchema.IsComplex() {

	}
	return nil
}
