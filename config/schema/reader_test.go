package schema

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs"
	"github.com/viant/assertly"
	"io"
	"io/ioutil"
	"smirror/config"
	"strings"
	"testing"
)

func TestTransformer_Read(t *testing.T) {

		maxBadRecord := 2
	var useCases = []struct {
		description string
		hasError    bool
		rule        *config.Rule
		reader      io.Reader
		expect      string
	}{

		{

			description: "csv validation  reader",
			rule: &config.Rule{
				Schema: &config.Schema{
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

			description: "quoted csv validation  reader",
			rule: &config.Rule{
				Schema: &config.Schema{
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
2,name2,"d,e,f",10`,
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

			description: "json validation  reader",

			rule: &config.Rule{
				Schema: &config.Schema{
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


		{

			description: "json validation  reader",

			rule: &config.Rule{
				Schema: &config.Schema{
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


		{

			description: "schema validation",

			rule: &config.Rule{
				Schema: &config.Schema{
					Format:        "JSON",
					MaxBadRecords: &maxBadRecord,
					Fields: []*config.Field{
						{
							Name:     "id",
							DataType: "int",
						},
						{
							Name:     "name",
							DataType: "string",
						},
						{
							Name:"Timestamp",
							DataType: "time",
							SourceDateFormat: "yyyy-MM-dd hh:mm:ss",
						},
					},
				},
			},
			hasError: false,
			reader: strings.NewReader(`{"id":1, "name":"name1", "Timestamp": "2019-07-25 01:03:22"}
{"id":"a", "name":"name7" }
{"id":3, "name":"name3"}`),
			expect: `{"id":1, "name":"name1", "Timestamp": "2019-07-25 01:03:22"}
{"id":3, "name":"name3" }
`,
		},

		{

			description: "csv schema ",
			rule: &config.Rule{
				Schema: &config.Schema{
					Format:     "CSV",
					Delimiter:  ",",
					Fields: []*config.Field{
						{
							Name:     "id",
							DataType: "int",
							Position:intPointer(0),
						},
						{
							Name:     "name",
							DataType: "string",
							Position:intPointer(1),
						},
						{
							Name:"Timestamp",
							DataType: "time",
							SourceDateFormat: "yyyy-MM-dd hh:mm:ss",
							Position:intPointer(2),

						},
						{
							Name:     "segment",
							DataType: "int",
							Position:intPointer(3),
						},
					},
				},
			},
			reader: strings.NewReader(`1,"event 1","2019-07-25 01:03:22", 3
2,"event 2","201j", 3
3,"event 3","2019-07-25 01:03:22",3.1
`),
			expect: `1,event 1,2019-07-25 01:03:22,3
3,event 3,2019-07-25 01:03:22,3`,
		},

	}

	for _, useCase := range useCases {
		_ = useCase.rule.Init(context.Background() ,afs.New())
		transformer, err := NewReader(useCase.reader, useCase.rule)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		data, err := ioutil.ReadAll(transformer)
		if useCase.hasError {
			assert.NotNil(t, err, useCase.description)
			continue
		}
		assert.Nil(t, err, useCase.description)
		assertly.AssertValues(t, useCase.expect, string(data))
	}

}


func intPointer(i int) *int {
	return &i
}