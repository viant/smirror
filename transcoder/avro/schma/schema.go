package schma

import (
	"encoding/json"
	"fmt"
	"github.com/viant/toolbox"
	"io"
	"strings"
)

//Schema represents a schema type
type Schema struct {
	Name       string                 `json:"name,omitempty"`
	Namespace  string                 `json:"namespace,omitempty"`
	Aliases    []string               `json:"aliases,omitempty"`
	Doc        string                 `json:"doc,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Default    interface{}            `json:"default,omitempty"`
	*Base
	*Union
	*Record
	*Array
	*Enum
	*Fixed

	translator Translator
}

func (t *Schema) SetTranslator(translator Translator) {
	t.translator = translator
}

//Translator writes value to writer using avro binary format
func (t *Schema) Write(value interface{}, w io.Writer) error {
	if t.translator == nil {
		return fmt.Errorf("translator was empty: %v", t.Type)
	}
	return t.translator(value, w)
}

func New(schemaValue interface{}) (*Schema, error) {
	result, err := newSchema(schemaValue)
	if err != nil {
		return nil, err
	}
	if result.Type == "" {
		return nil, fmt.Errorf("type was empty for %s", schemaValue)
	}
	return result, err
}

func newSchema(schemaValue interface{}) (*Schema, error) {
	elemType := &Schema{Base: &Base{}}
	elemType.Type = toolbox.AsString(schemaValue)
	switch val := schemaValue.(type) {
	case toolbox.AnyJSONType:
		return New([]byte(val))
	case []byte:
		if json.Valid(val) {
			JSON := strings.TrimSpace(string(val))

			if strings.HasPrefix(JSON, `"`) || JSON == "null" {
				elemType.Type = strings.Trim(string(val), `"`)
				return elemType, nil
			}
			if strings.HasPrefix(JSON, "[") {
				array := []interface{}{}
				if err := json.Unmarshal(val, &array); err != nil {
					return nil, err
				}
				return New(array)
			}
			aMap := map[string]interface{}{}
			if err := json.Unmarshal(val, &aMap); err != nil {
				return nil, err
			}
			return New(aMap)
		}

	case string:
		return New([]byte(val))
	case map[string]interface{}:
		JSON, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(JSON, &elemType)
		if err != nil {
			return nil, err
		}

		if elemType.IsRecord() && len(elemType.Fields) > 0 {
			for i := range elemType.Fields {
				elemType.Fields[i].Type, err = New(elemType.Fields[i].AnyType)
				if err != nil {
					return nil, err
				}
				elemType.Fields[i].AnyType = ""
			}
		}
		if elemType.IsArray() {
			if elemType.Array.Items, err = New(elemType.AnyItems); err != nil {
				return elemType, err
			}
			elemType.AnyItems = ""
		}
		return elemType, err
	case []interface{}:
		elemType.Type = typeUnion
		if elemType.Union == nil {
			elemType.Union = &Union{}
		}
		elemType.Types = make([]*Schema, 0)

		for i := range val {
			unionType, err := New(val[i])
			if err != nil {
				return nil, err
			}
			elemType.Union.Types = append(elemType.Union.Types, unionType)
		}
		return elemType, nil

	default:
		return nil, fmt.Errorf("unsupported type: %T %v", schemaValue, schemaValue)
	}

	return elemType, nil
}
