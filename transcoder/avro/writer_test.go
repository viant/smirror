package avro

import (
	"bytes"
	"github.com/actgardner/gogen-avro/container"
	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
	"github.com/viant/toolbox"
	"github.com/viant/smirror/transcoder/avro/schma"
	"testing"
)

func TestSetWriter(t *testing.T) {

	var useCases = []struct {
		description string
		schema      string
		data        map[string]interface{}
		expect      string
	}{
		{
			description: "string type",
			schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "foo",
		"fields": [
			{ "name": "bar", "type": "string"}
		]
	}`,
			data: map[string]interface{}{
				"bar": "test 1",
			},
		},
		{
			description: "int type",
			schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "foo",
		"fields": [
			{ "name": "bar", "type": "int"}
		]
	}`,
			data: map[string]interface{}{
				"bar": 11,
			},
		},
		{
			description: "boolean type",
			schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "foo",
		"fields": [
			{ "name": "bar", "type": "boolean"}
		]
	}`,
			data: map[string]interface{}{
				"bar": true,
			},
		},
		{
			description: "double type",
			schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "foo",
		"fields": [
			{ "name": "bar", "type": "double"}
		]
	}`,
			data: map[string]interface{}{
				"bar": 3.3,
			},
		},
		{
			description: "float type",
			schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "foo",
		"fields": [
			{ "name": "bar", "type": "float"}
		]
	}`,
			data: map[string]interface{}{
				"bar": 3.3,
			},
		},

		{
			description: "union with primitive type and null",
			schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "foo",
		"fields": [
			{ "name": "id", "type": "int"},
			{ "name": "name", "type": ["null", "string"]}
		]
	}`,
			data: map[string]interface{}{
				"id": 1,
			},
		},

		{
			description: "union with array - non empty ",
			schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "foo",
		"fields": [
			{ "name": "id", "type": "int"},
			{ "name": "errors", "type": ["null", {"type":"array", "items":"string"}], "default": null }
		]
	}`,
			data: map[string]interface{}{
				"id":     1,
				"errors": []interface{}{"error 1", "error 2"},
			},
		},
		{
			description: "union with array - empty ",
			schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "foo",
		"fields": [
			{ "name": "id", "type": "int"},
			{ "name": "errors", "type": ["null", {"type":"array", "items":"int"}], "default": null }
		]
	}`,
			data: map[string]interface{}{
				"id":     1,
				"errors": []interface{}{},
			},
		},
		{
			description: "union with record ",
			schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "foo",
		"fields": [
			{ "name": "id", "type": "int"},
			{ "name": "foo", "type": ["null",{
				"type":	"record",
				"name": "foo",
				"fields": [
					{ "name": "id", "type": "int"}
				]
			}],"default":null}
		]
	}`,
			data: map[string]interface{}{
				"id": 1,
				"foo": map[string]interface{}{
					"id": 2,
				},
			},
		},
		{
			description: "union with repeated record ",
			schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "foo",
		"fields": [
			{ "name": "id", "type": "int"},
			{ "name": "foo", "type": ["null", {"type":"array", "items":{
				"type":	"record",
				"name": "foo",
				"fields": [
					{ "name": "id", "type": "int"}
				]
			}}], "default": null }
		]
	}`,
			data: map[string]interface{}{
				"id": 1,
				"foo": []interface{}{
					map[string]interface{}{
						"id": 2,
					},
				},
			},
		},
	}

	for _, useCase := range useCases {

		schema, err := schma.New(useCase.schema)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		err = SetWriter(schema)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		buffer := new(bytes.Buffer)
		writer, err := container.NewWriter(buffer, container.Snappy, 40, useCase.schema)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		record := NewRecord(useCase.data, schema, useCase.schema)
		err = writer.WriteRecord(record)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		err = writer.Flush()
		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		avroReader, err := goavro.NewOCFReader(buffer)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		assert.True(t, avroReader.Scan())
		actual, err := avroReader.Read()
		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		actualMap := toolbox.AsMap(actual)
		transformUnions(schema, actualMap)

		if !assertly.AssertValues(t, useCase.data, actualMap, useCase.description) {
			_ = toolbox.DumpIndent(actual, true)
		}
	}

}

func transformUnions(schema *schma.Schema, actualMap map[string]interface{}) {
	for _, field := range schema.Fields {
		if !field.Type.IsUnion() {
			continue
		}
		unionValue, ok := actualMap[field.Name]
		if !ok {
			continue
		}
		delete(actualMap, field.Name)
		unionMap := toolbox.AsMap(unionValue)
		for _, v := range unionMap {
			actualMap[field.Name] = v
		}
	}
}
