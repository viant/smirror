package avro

import (
	"fmt"
	"github.com/viant/toolbox"
	"io"
	"smirror/transcoder/avro/schma"
	"strings"
	"time"
)

func translateToLogicalTime(schema *schma.Schema) schma.Translator {
	return func(value interface{}, w io.Writer) error {
		if toolbox.IsInt(value) {
			intValue := toolbox.AsInt(value)
			return writeLong(int64(intValue), w)
		}
		var err error
		var ts *time.Time
		switch val := value.(type) {
		case string:
			layout := defaultTimeLayout
			if strings.Contains(val, "T") {
				layout = time.RFC3339Nano
			}
			if ts, err = toolbox.ToTime(val, layout); err != nil {
				return err
			}
		case *time.Time:
			ts = val
		case time.Time:
			ts = &val
		default:
			return fmt.Errorf("unsupported logical time type: %T", value)
		}
		var numericTime int64
		if strings.Contains(schema.LogicalType, millis) { //nano to millis
			numericTime = ts.UnixNano() / 1000000
		} else { //nano to micros
			numericTime = ts.UnixNano() * 1000
		}
		return writeLong(numericTime, w)
	}
}
