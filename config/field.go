package config

import (
	"github.com/viant/toolbox"
	"strings"
)

const (
	DataTypeTime = "time"
	DataTypeFloat = "float"
	DataTypeInt = "int"
	DataTypeBoolean = "boolean"
	DataTypeString = "string"
)

//Fields represents recovery filed
type Field struct {
	Name string
	Position *int
	DataType string
	SourceDateFormat string
	SourceDateLayout string
	TargetDateFormat string
	TargetDateLayout string
}

func (f *Field) Init() {
	if f.SourceDateLayout == "" && f.SourceDateFormat != "" {
		f.SourceDateLayout = toolbox.DateFormatToLayout(f.SourceDateFormat)
	}
	if f.TargetDateLayout == "" && f.TargetDateFormat != "" {
		f.TargetDateLayout = toolbox.DateFormatToLayout(f.TargetDateFormat)
	}
}



//AdjustText adjust text value
func (f *Field) AdjustValue(value interface{}) interface{} {
	switch strings.ToLower(f.DataType) {
	case DataTypeTime:
		source , err:= toolbox.ToTime(value, f.SourceDateLayout)
		if err != nil {
			return value
		}
		return source.Format(f.TargetDateLayout)
	case DataTypeFloat:
		result, err := toolbox.ToFloat(value)
		if err != nil {
			return value
		}
		return result
	case DataTypeInt:
		result, err := toolbox.ToInt(value)
		if err != nil {
			return value
		}
		return result
	case DataTypeBoolean:
		result, err := toolbox.ToBoolean(value)
		if err != nil {
			return value
		}
		return result
	}
	return value

}

//AdjustText adjust text value
func (f *Field) AdjustText(value string) string {
	switch strings.ToLower(f.DataType) {
	case DataTypeTime:
		source , err:= toolbox.ToTime(value, f.SourceDateLayout)
		if err != nil {
			return value
		}
		return source.Format(f.TargetDateLayout)
	case DataTypeFloat:
		result, err := toolbox.ToFloat(value)
		if err != nil {
			return value
		}
		return toolbox.AsString(result)
	case DataTypeInt:
		result, err := toolbox.ToInt(value)
		if err != nil {
			return value
		}
		return toolbox.AsString(result)
	case DataTypeBoolean:
		result, err := toolbox.ToBoolean(value)
		if err != nil {
			return value
		}
		return toolbox.AsString(result)
	}
	return value

}