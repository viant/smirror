package schma

import (
	"github.com/viant/toolbox"
)

type Field struct {
	Name       string                 `json:"name,omitempty"`
	Doc        string                 `json:"doc,omitempty"`
	Default    interface{}            `json:"default"`
	AnyType    toolbox.AnyJSONType    `json:"type,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Type       *Schema
}
