package avro

import (
	"fmt"
	"smirror/transcoder/avro/schma"
)

//SetWriter initialises avro schema
func SetWriter(schema *schma.Schema) error {
	if schema.IsComplex() {
		if schema.IsRecord() {
			schema.SetTranslator(translateToRecord(schema))
			for i := range schema.Fields {
				if err := SetWriter(schema.Fields[i].Type); err != nil {
					return err
				}
			}
			return nil
		}
		if schema.IsUnion() {
			schema.SetTranslator(translateToUnion(schema))
			for i := range schema.Types {
				err := SetWriter(schema.Types[i])
				if err != nil {
					return err
				}
			}
			return nil

		}
		if schema.IsArray() {
			schema.SetTranslator(translateToArray(schema.Items))
			if err := SetWriter(schema.Items); err != nil {
				return err
			}
		}
		return nil
	}

	switch schema.LogicalType {
	//case "date": // 	DATE
	case timeMicros, timeMillis, timestampMicros, timestampMillis:
		schema.SetTranslator(translateToLogicalTime(schema))
		return nil
	}
	if schema.IsString() || schema.IsString() {
		schema.SetTranslator(translateToString)
		return nil
	}

	if schema.IsInt() || schema.IsLong() {
		schema.SetTranslator(translateToLong)
		return nil
	}
	if schema.IsDouble() {
		schema.SetTranslator(translateToDouble)
		return nil
	}
	if schema.IsFloat() {
		schema.SetTranslator(translateToFloat)
		return nil
	}
	if schema.IsBoolean() {
		schema.SetTranslator(translateToBoolean)
		return nil
	}
	if schema.IsNull() {
		schema.SetTranslator(translateToNull)
		return nil
	}
	if schema.IsMap() {
		return fmt.Errorf("not supported yet: %v", schema.Type)
	}
	return fmt.Errorf("failed to set translator for: %v", schema.Type)
}
