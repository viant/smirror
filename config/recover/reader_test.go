package recover

import (
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"smirror/config"
	"strings"
	"testing"
)

func TestTransformer_Read(t *testing.T) {

	var useCases = []struct {
		description string
		rule        *config.Rule
		reader      io.Reader
		expect      string
	}{

		{

			description: "csv recovery reader",
			rule: &config.Rule{
				Recover: &config.Recover{
					Format:     "CSV",
					Delimiter:  ",",
					FieldCount: 3,
				},
			},
			reader: strings.NewReader(`
1,name1,desc1
2,name2,desc2,ee
3,name3
4,name4
`),
			expect: `1,name1,desc1
2,name2,desc2
3,name3,
4,name4,`,
		},

		{

			description: "quoted csv recovery reader",
			rule: &config.Rule{
				Recover: &config.Recover{
					Format:     "CSV",
					Delimiter:  ",",
					FieldCount: 4,
				},
			},
			reader: strings.NewReader(`
1,"name1","a,b,c"
2,"name2","d,e,f",10
`),
			expect: `1,name1,"a,b,c",
2,"name2","d,e,f",10`,
		},

		{

			description: "replacement reader",
			rule: &config.Rule{
				Replace: []*config.Replace{
					{
						From: "a",
						To:   "b",
					},
				},
			},
			reader: strings.NewReader("abcde\n  ghijka\n   zsa"),
			expect: "bbcde\n  ghijkb\n   zsb",
		},

		{

			description: "json recovery reader",

			rule: &config.Rule{
				Recover: &config.Recover{
					Format: "JSON",
				},
			},
			reader: strings.NewReader(`

{"id":1, "name":"name1"}
{"id":2, "name":"name2" 
{"id":3, "name":"name3"}

`),
			expect: `{"id":1, "name":"name1"}
{"id":3, "name":"name3"}`,
		},
	}

	for _, useCase := range useCases {

		transformer, err := NewReader(useCase.reader, useCase.rule)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		data, err := ioutil.ReadAll(transformer)
		assert.Nil(t, err, useCase.description)
		assert.EqualValues(t, useCase.expect, string(data))

	}

}
