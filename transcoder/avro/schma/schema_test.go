package schma

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
	"github.com/viant/toolbox"
	"testing"
)

func TestInitSchema(t *testing.T) {

	var useCases = []struct {
		description string
		schema      string
		expect      string
	}{

		{
			description: "nested",
			schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "indentity",
		"fields": [
			{ "name": "FirstName", "type": "string"},
			{ "name": "LastName", "type": "string"},
			{ "name": "Errors", "type": ["null", {"type":"array", "items":"string"}], "default": null },
			{ "name": "Address", "type": ["null",{
				"namespace": "my.namespace.com",
				"type":	"record",
				"name": "address",
				"fields": [
					{ "name": "Address1", "type": "string" },
					{ "name": "Address2", "type": ["null", "string"], "default": null },
					{ "name": "City", "type": "string" },
					{ "name": "State", "type": "string" },
					{ "name": "Zip", "type": "int" }
				]
			}],"default":null}
		]
	}`,

			expect: `{
	"Fields": [
		{
			"Name": "FirstName",
			"Type": {
				"Type": "string"
			}
		},
		{
			"Name": "LastName",
			"Type": {
				"Type": "string"
			}
		},
		{
			"Name": "Errors",
			"Type": {
				"Type": "union",
				"Types": [
					{
						"Type": "null"
					},
					{
						"Items": {
							"Type": "string"
						},
						"Type": "array"
					}
				]
			}
		},
		{
			"Name": "Address",
			"Type": {
				"Type": "union",
				"Types": [
					{
						"Type": "null"
					},
					{
						"Fields": [
							{
								"Name": "Address1",
								"Type": {
									"Type": "string"
								}
							},
							{
								"Name": "Address2",
								"Type": {
									"Type": "union"
								}
							},
							{
								"Name": "City",
								"Type": {
									"Type": "string"
								}
							},
							{
								"Name": "State",
								"Type": {
									"Type": "string"
								}
							},
							{
								"Name": "Zip",
								"Type": {
									"Type": "int"
								}
							}
						],
						"Name": "address",
						"Namespace": "my.namespace.com",
						"Type": "record"
					}
				]
			}
		}
	],
	"Name": "indentity",
	"Namespace": "my.namespace.com",
	"Type": "record"
}`,
		},
	}

	for _, useCase := range useCases {
		schema, err := New(useCase.schema)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		JSON, _ := json.MarshalIndent(schema, "\t", " ")
		fmt.Printf("%s\n", JSON)
		if !assertly.AssertValues(t, useCase.expect, schema) {
			toolbox.DumpIndent(schema, true)
		}
	}

}
