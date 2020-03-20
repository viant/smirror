package config

import (
	"github.com/viant/toolbox"
	"strings"
)

const (
	DataTypeTime    = "time"
	DataTypeFloat   = "float"
	DataTypeInt     = "int"
	DataTypeBoolean = "boolean"
	DataTypeString  = "string"
)

//Fields represents validation filed
type Field struct {
	Name             string
	Position         *int
	DataType         string
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
func (f *Field) AdjustValue(value interface{}) (interface{}, error) {
	switch strings.ToLower(f.DataType) {
	case DataTypeTime:
		if f.SourceDateLayout == "" && f.TargetDateLayout == "" {
			return value, nil
		}
		if f.SourceDateLayout == "" && f.TargetDateLayout != "" {
			f.SourceDateLayout = f.TargetDateLayout
		}
		source, err := toolbox.ToTime(value, f.SourceDateLayout)
		if err != nil {
			return value, err
		}
		if f.TargetDateLayout == "" {
			f.TargetDateLayout = f.SourceDateLayout
		}
		return source.Format(f.TargetDateLayout), nil
	case DataTypeFloat:
		result, err := toolbox.ToFloat(value)
		if err != nil {
			return value, err
		}
		return result, nil
	case DataTypeInt:
		result, err := toolbox.ToInt(value)
		if err != nil {
			return value, err
		}
		return result, nil
	case DataTypeBoolean:
		result, err := toolbox.ToBoolean(value)
		if err != nil {
			return value, err
		}
		return result, err
	}
	return value, nil

}

//AdjustText adjust text value
func (f *Field) AdjustText(value string) (string, error) {

	switch strings.ToLower(f.DataType) {
	case DataTypeTime:
		value = strings.TrimSpace(value)
		if f.SourceDateLayout == "" && f.TargetDateLayout == "" {
			return value, nil
		}
		if f.TargetDateLayout != "" && f.SourceDateLayout == "" {
			f.SourceDateLayout = f.TargetDateLayout
		}
		source, err := toolbox.ToTime(value, f.SourceDateLayout)
		if err != nil {
			return "", err
		}
		if f.TargetDateLayout == "" {
			f.TargetDateLayout = f.SourceDateLayout
		}
		return source.Format(f.TargetDateLayout), nil
	case DataTypeFloat:
		value = strings.TrimSpace(value)
		result, err := toolbox.ToFloat(value)
		if err != nil {
			return value, err
		}
		return toolbox.AsString(result), nil
	case DataTypeInt:
		value = strings.TrimSpace(value)
		result, err := toolbox.ToInt(value)
		if err != nil {
			return value, err
		}
		return toolbox.AsString(result), nil
	case DataTypeBoolean:
		value = strings.TrimSpace(value)
		result, err := toolbox.ToBoolean(value)
		if err != nil {
			return value, err
		}
		return toolbox.AsString(result), nil
	}
	return value, nil
}
