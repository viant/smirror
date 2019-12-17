package avro

import (
	"github.com/viant/toolbox"
	"io"
	"smirror/transcoder/avro/schma"
	"strings"
)

func translateToArray(itemSchema *schma.Schema) schma.Translator {
	return func(value interface{}, w io.Writer) error {
		aSlice := convertToSlice(value)
		err := writeLong(int64(len(aSlice)), w)
		if err != nil || len(aSlice) == 0 {
			return err
		}
		for i := range aSlice {
			if err = itemSchema.Write(aSlice[i], w); err != nil {
				return err
			}
		}
		return writeLong(0, w)
	}
}

func convertToSlice(value interface{}) []interface{} {
	switch value.(type) {
	case string, []byte:
		text := toolbox.AsString(value)
		if strings.Count(text, " ") > 0 {
			value = strings.Split(text, " ")
		} else if strings.Count(text, ",") > 0 {
			value = strings.Split(text, ",")
		}
	}
	return toolbox.AsSlice(value)
}
