package transcoder

import (
	"bytes"
	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
	"github.com/viant/toolbox"
	"io/ioutil"
	"smirror/config"
	"smirror/config/transcoding"
	"smirror/transcoder/avro/schma"
	"strings"
	"testing"
)

func TestReader_Read(t *testing.T) {
	var useCases = []struct {
		description string
		input       string
		config.Transcoding
		expect interface{}
	}{
		{
			description: "CSV to avro",

			input: `1,name 1,desc 1
2,name 2,desc 2,
3,name 3,desc 3`,

			expect: []map[string]interface{}{
				{
					"id":          1,
					"name":        "name 1",
					"description": "desc 1",
				},
				{
					"id":          2,
					"name":        "name 2",
					"description": "desc 2",
				},
				{
					"id":          3,
					"name":        "name 3",
					"description": "desc 3",
				},
			},
			Transcoding: config.Transcoding{
				Source: transcoding.Codec{
					Format: "CSV",
					Fields: []string{"id", "name", "description"},
				},
				Dest: transcoding.Codec{
					Format: "AVRO",
					Schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "foo",
		"fields": [
			{ "name": "id", "type": "int"},
		    { "name": "name", "type": "string"},
			{ "name": "description", "type": "string"}
		]
	}`,
				},
			},
		},
		{
			description: "JSON to avro",

			input: `{"id":1,"name":"name 1","description":"desc 1"}
{"id":2,"name":"name 2","description":"desc 2"}
{"id":3,"name":"name 3","description":"desc 3"}`,

			expect: []map[string]interface{}{
				{
					"id":          1,
					"name":        "name 1",
					"description": "desc 1",
				},
				{
					"id":          2,
					"name":        "name 2",
					"description": "desc 2",
				},
				{
					"id":          3,
					"name":        "name 3",
					"description": "desc 3",
				},
			},
			Transcoding: config.Transcoding{
				Source: transcoding.Codec{
					Format: "JSON",
				},
				Dest: transcoding.Codec{
					Format: "AVRO",
					Schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "foo",
		"fields": [
			{ "name": "id", "type": "int"},
		    { "name": "name", "type": "string"},
			{ "name": "description", "type": "string"}
		]
	}`,
				},
			},
		},
		{
			description: "CSV to avro with mapping",

			input: `1,name 1,desc 1
2,name 2,desc 2,
3,name 3,desc 3`,

			expect: []map[string]interface{}{
				{
					"id": 1,
					"attr1": map[string]interface{}{
						"name":        "name 1",
						"description": "desc 1",
					},
				},
				{
					"id": 2,
					"attr1": map[string]interface{}{
						"name":        "name 2",
						"description": "desc 2",
					},
				},
				{
					"id": 3,
					"attr1": map[string]interface{}{
						"name":        "name 3",
						"description": "desc 3",
					},
				},
			},
			Transcoding: config.Transcoding{
				Source: transcoding.Codec{
					Format: "CSV",
					Fields: []string{"id", "name", "description"},
				},
				PathMapping: transcoding.Mappings{
					{
						From: "id",
						To:   "id",
					},
					{
						From: "name",
						To:   "attr1.name",
					},

					{
						From: "description",
						To:   "attr1.description",
					},
				},
				Dest: transcoding.Codec{
					Format: "AVRO",
					Schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "root",
		"fields": [
			{ "name": "id", "type": "int"},
			{ "name": "attr1", "type": ["null",{
				"type":	"record",
				"name": "foo",
				"fields": [
					{ "name": "name", "type": "string"},
					{ "name": "description", "type": "string"}
				]
			}],"default":null}
		]
	}`,
				},
			},
		},
	}

	for _, useCase := range useCases {

		reader, err := NewReader(strings.NewReader(useCase.input), &useCase.Transcoding)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		schema, err := schma.New(useCase.Transcoding.Dest.Schema)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		data, err := ioutil.ReadAll(reader)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		avroReader, err := goavro.NewOCFReader(bytes.NewReader(data))
		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		var actual = make([]interface{}, 0)
		for avroReader.Scan() {

			actualRecords, err := avroReader.Read()
			if !assert.Nil(t, err, useCase.description) {
				continue
			}
			actualMap := toolbox.AsMap(actualRecords)
			transformUnions(schema, actualMap)
			actual = append(actual, actualMap)
		}

		if !assertly.AssertValues(t, useCase.expect, actual, useCase.description) {
			toolbox.DumpIndent(actual, true)
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
