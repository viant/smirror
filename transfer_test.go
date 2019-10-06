package smirror

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"smirror/config"
	"strings"
	"testing"
)

func TestTransfer_GetReader(t *testing.T) {

	var useCases = []struct {
		description string
		data        string
		replace     map[string]string
		expect      string
	}{
		{
			description: "basic replacement",
			data: `this is my line 1
this is my line 2
`,
			replace: map[string]string{
				"is": "abc",
			},
			expect: `thabc abc my line 1
thabc abc my line 2`,
		},
		{
			description: "multi line replacement",
			data: `
"123213"",""
adaew"",""
""1`,
			replace: map[string]string{
				`""`: `"`,
			},
			expect: `
"123213","
adaew","
"1`,
		},
	}

	for _, useCase := range useCases {

		var replace = make([]*config.Replace, 0)
		for k, v := range useCase.replace {
			replace = append(replace, &config.Replace{From: k, To: v})
		}
		transfer := &Transfer{

			rewriter: NewRewriter(replace...),
			Reader:   strings.NewReader(useCase.data),
		}
		reader, err := transfer.GetReader()
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		actual, err := ioutil.ReadAll(reader)
		assert.Nil(t, err, useCase.description)
		assert.EqualValues(t, useCase.expect, string(actual), useCase.description)
	}

}
