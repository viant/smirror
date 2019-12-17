package avro

import (
	"github.com/viant/toolbox"
	"io"
	"smirror/transcoder/avro/schma"
)

func translateToUnion(schema *schma.Schema) schma.Translator {
	return func(value interface{}, w io.Writer) (err error) {
		if isValueNull(value) && schema.IsNullUnion() {
			_, err = writeUnionIndex(schema, true, w)
			return err
		}
		var index = 0
		if index, err = writeUnionIndex(schema, false, w); err != nil {
			return err
		}
		return schema.Types[index].Write(value, w)
	}
}

func writeUnionIndex(schema *schma.Schema, isNull bool, w io.Writer) (int, error) {
	for i, uType := range schema.Types {
		if isNull {
			if !uType.IsNull() {
				continue
			}
			return i, writeLong(int64(i), w)
		}
		if uType.IsNull() {
			continue
		}
		return i, writeLong(int64(i), w)
	}

	return -1, nil
}

func isValueNull(value interface{}) bool {
	if value == nil {
		return true
	}

	if toolbox.IsString(value) {
		return toolbox.AsString(value) == ""
	}

	if toolbox.IsFloat(value) {
		return toolbox.AsFloat(value) == 0
	}
	if toolbox.IsInt(value) {
		return toolbox.AsInt(value) == 0
	}
	if toolbox.IsBool(value) {
		return !toolbox.AsBoolean(value)
	}
	return false
}
