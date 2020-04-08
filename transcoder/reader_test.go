package transcoder

import (
	"bytes"
	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
	"github.com/viant/toolbox"
	"io/ioutil"
	"path"
	"smirror/config"
	"smirror/config/transcoding"
	"smirror/transcoder/avro/schma"
	"smirror/transcoder/xlsx"
	"testing"
)

func TestReader_Read(t *testing.T) {

	baseDir := toolbox.CallerDirectory(3)
	xlsData, err := ioutil.ReadFile(path.Join(baseDir, "test", "book.xlsx"))
	if ! assert.Nil(t, err) {
		return
	}
	
	var useCases = []struct {
		description string
		 isXls bool
		input       []byte
		config.Transcoding
		expect interface{}
	}{


		{
			isXls:true,
			description:"XLS to AVRO",
			input:xlsData,
			Transcoding: config.Transcoding{
				Source: transcoding.Codec{
					Format: "XLSX",
				},
				Dest: transcoding.Codec{
					Format:"AVRO",
					RecordPerBlock: 20,
				},
				Autodetect:true,
			},
			expect: []map[string]interface{}{
				{
					"Active": true,
					"Count": 1,
					"Id": 1,
					"Name": "Name 1",
					"Timestamp": 1550102400000,
					"Value": 1.2000000476837158,
				},
				{
					"Count": 3,
					"Id": 2,
					"Name": "Name 2",
					"Value": 4.3,
				},
				{
					"Count": 4,
					"Id": 3,
					"Name": "Name 3",
					"Timestamp": 1550102400000,
				},
			},
		},
		
		{
			description: "CSV to AVRO",

			input:[]byte( `1,name 1,desc 1
2,name 2,desc 2,
3,name 3,desc 3`),

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
					RecordPerBlock: 2,
					Format:         "AVRO",
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
			description: "JSON to AVRO",

			input:[]byte( `{"id":1,"name":"name 1","description":"desc 1"}
{"id":2,"name":"name 2","description":"desc 2"}
{"id":3,"name":"name 3","description":"desc 3"}`),

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
			description: "CSV to AVRO with mapping",

			input: []byte(`1,name 1,desc 1
2,name 2,desc 2,
3,name 3,desc 3`),

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

		{
			description: "tsv log to AVRO",

			input: []byte(`2019-12-16T23:55:38.199597Z ci-ad-vpc-east 107.77.216.9:55151 10.55.8.61:8080 0.000027 0.00059 0.000025 200 200 0 631 "GET https://tabc.comm:443/d/rt/pixel?rtsite_id=23762&uuid=8cd4f7f5-6b0a-4697-87eb-db4556736c57&rr=1936727420 HTTP/1.1" "Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.4 Mobile/15E148 Safari/604.1" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2
2019-12-16T23:55:38.206419Z ci-ad-vpc-east 47.212.89.174:52566 10.55.8.44:8080 0.000026 0.000463 0.000026 200 200 0 0 "GET https://tabc.com.com:443/d/track/video?zid=googleadx_1_0_1&sid=87ddc05f-205f-11ea-a9c6-55bf429b7234&crid=19841312&adid=54517&oid=1114112&cid=172922&spid=282&pubid=45&site_id=442320&auid=1452372&algid=0&algrev=0&offpc=0&maxbid=0.000&optpc=0&cstpc=0&ez_p=&eid=2 HTTP/1.1" "Mozilla/5.0 (iPhone; CPU iPhone OS 12_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 [FBAN/FBIOS;FBDV/iPhone10,2;FBMD/iPhone;FBSN/iOS;FBSV/12.4.1;FBSS/3;FBID/phone;FBLC/en_US;FBOP/5;FBCR/Verizon]" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2`),

			expect: []map[string]interface{}{
				{
					"backend_status_code": 200,
					"request":             "GET https://tabc.comm:443/d/rt/pixel?rtsite_id=23762\u0026uuid=8cd4f7f5-6b0a-4697-87eb-db4556736c57\u0026rr=1936727420 HTTP/1.1",
					"timestamp":           1576540538199,
					"user_agent":          "Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.4 Mobile/15E148 Safari/604.1",
				},
				{
					"backend_status_code": 200,
					"request":             "GET https://tabc.com.com:443/d/track/video?zid=googleadx_1_0_1\u0026sid=87ddc05f-205f-11ea-a9c6-55bf429b7234\u0026crid=19841312\u0026adid=54517\u0026oid=1114112\u0026cid=172922\u0026spid=282\u0026pubid=45\u0026site_id=442320\u0026auid=1452372\u0026algid=0\u0026algrev=0\u0026offpc=0\u0026maxbid=0.000\u0026optpc=0\u0026cstpc=0\u0026ez_p=\u0026eid=2 HTTP/1.1",
					"timestamp":           1576540538206,
					"user_agent":          "Mozilla/5.0 (iPhone; CPU iPhone OS 12_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 [FBAN/FBIOS;FBDV/iPhone10,2;FBMD/iPhone;FBSN/iOS;FBSV/12.4.1;FBSS/3;FBID/phone;FBLC/en_US;FBOP/5;FBCR/Verizon]",
				},
			},
			Transcoding: config.Transcoding{
				Source: transcoding.Codec{
					Format:    "CSV",
					Delimiter: " ",
					Fields: []string{"timestamp", "elb", "client_port", "backend_port", "request_processing_time", "backend_processing_time", "response_processing_time", "elb_status_code", "backend_status_code",
						"received_bytes", "sent_bytes", "request", "user_agent", "ssl_cipher", "ssl_protocol",
					},
				},
				Dest: transcoding.Codec{
					RecordPerBlock: 1,
					Format:         "AVRO",
					Schema: `{
		"namespace": "my.namespace.com",
		"type":	"record",
		"name": "foo",
		"fields": [
			{ "name": "timestamp", "type" : {"type" : "long", "logicalType" : "timestamp-millis"}},
		    { "name": "backend_status_code", "type": "int"},
			{ "name": "request", "type": "string"},
			{ "name": "user_agent", "type": "string"}
		]
	}`,
				},
			},
		},
	}

	for _, useCase := range useCases[:1] {

		reader, err := NewReader(bytes.NewReader(useCase.input), &useCase.Transcoding, 0)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}


		rawSchema := useCase.Transcoding.Dest.Schema
		if useCase.isXls {
			decoder, _ := xlsx.NewDecoder(bytes.NewReader(useCase.input))
			rawSchema = decoder.Schema()
		}

		schema, err := schma.New(rawSchema)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		data, err := ioutil.ReadAll(reader)
		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		AVROReader, err := goavro.NewOCFReader(bytes.NewReader(data))

		if !assert.Nil(t, err, useCase.description) {
			continue
		}

		assert.Nil(t, AVROReader.Err(), useCase.description)
		var actual = make([]interface{}, 0)
		for AVROReader.Scan() {

			actualRecords, err := AVROReader.Read()
			if !assert.Nil(t, err, useCase.description) {
				continue
			}
			actualMap := toolbox.AsMap(actualRecords)
			transformUnions(schema, actualMap)
			actual = append(actual, actualMap)
		}

		assert.Nil(t, AVROReader.Err(), useCase.description)

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
