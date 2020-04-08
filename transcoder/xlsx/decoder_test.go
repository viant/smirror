package xlsx

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
	"github.com/viant/toolbox"
	"io/ioutil"
	"path"
	"testing"
)

func TestDecoder_Schema(t *testing.T) {

	baseDir := toolbox.CallerDirectory(3)

	var useCases = []struct{
		description string
		location string
		expect string
	} {
		{
			description:"basic schema",
			location:path.Join(baseDir ,"test", "book.xlsx"),
			expect:`{"type":	"record",
"name": "Root",
"fields": [{ "name": "Id", "type": ["null", "long"], "default": null},
{ "name": "Name", "type": ["null", "string"], "default": null},
{ "name": "Value", "type": ["null", "float"], "default": null},
{ "name": "Count", "type": ["null", "string"], "default": null},
{ "name": "Timestamp", "type" : ["null", {"type" : "long", "logicalType" : "timestamp-millis"}], "default": null},
{ "name": "Active", "type": ["null", "boolean"], "default": null}]
}`,
		},

	}

	for _, useCase := range useCases {
		data, err := ioutil.ReadFile(useCase.location)
		if ! assert.Nil(t, err, useCase.description) {
			continue
		}
		decoder, err := NewDecoder(bytes.NewReader(data))
		assert.Nil(t, err)
		actual := decoder.Schema()
		if ! assertly.AssertValues(t, useCase.expect, actual, useCase.description) {
			fmt.Printf("!!! %s\n", actual)
		}
	}
}
