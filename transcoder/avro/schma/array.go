package schma

import "github.com/viant/toolbox"

type Array struct {
	Items    *Schema
	AnyItems toolbox.AnyJSONType `json:"items,omitempty"`
}
