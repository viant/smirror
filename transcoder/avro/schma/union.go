package schma

import "github.com/viant/toolbox"

//Union represents a union type
type Union struct {
	Types    []*Schema
	AnyTypes toolbox.AnyJSONType `column:"types"`
}

//IsNullUnion returns true if any union type is a null type
func (u *Union) IsNullUnion() bool {
	for _, uType := range u.Types {
		if uType.IsNull() {
			return true
		}
	}
	return false
}
