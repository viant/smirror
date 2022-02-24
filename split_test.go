package smirror

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"github.com/viant/smirror/config"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {

	var data = make([]byte, 0)
	for i := 0; i < 9; i++ {
		data = append(data, []byte(fmt.Sprintf("%v %v\n", strings.Repeat("x", 2), i))...)
	}
	text := string(data)
	index := 0
	useCases := []struct {
		description string
		rule        *config.Rule

		expect []string
	}{
		{
			description: "no more then 4 lines case",
			rule: &config.Rule{
				Split: &config.Split{MaxLines: 4},
			},
			expect: []string{
				"xx 0\nxx 1\nxx 2\nxx 3", "xx 4\nxx 5\nxx 6\nxx 7", "xx 8",
			},
		},
		{
			description: "3 elements each",
			rule: &config.Rule{
				Split: &config.Split{MaxLines: 3},
			},
			expect: []string{
				"xx 0\nxx 1\nxx 2", "xx 3\nxx 4\nxx 5", "xx 6\nxx 7\nxx 8",
			},
		},
		{
			description: "9 elements",
			rule: &config.Rule{
				Split: &config.Split{MaxLines: 1},
			},
			expect: []string{
				"xx 0", "xx 1", "xx 2", "xx 3", "xx 4", "xx 5", "xx 6", "xx 7", "xx 8",
			},
		},
		{
			description: "9 elements",
			rule: &config.Rule{
				Split: &config.Split{MaxLines: 0},
			},
			expect: []string{
				"xx 0", "xx 1", "xx 2", "xx 3", "xx 4", "xx 5", "xx 6", "xx 7", "xx 8",
			},
		},
		{
			description: "1 elements",
			rule: &config.Rule{
				Split: &config.Split{MaxLines: 10},
			},
			expect: []string{
				"xx 0\nxx 1\nxx 2\nxx 3\nxx 4\nxx 5\nxx 6\nxx 7\nxx 8",
			},
		},
		{
			description: "1 elements",
			rule: &config.Rule{
				Split: &config.Split{MaxLines: 9},
			},
			expect: []string{
				"xx 0\nxx 1\nxx 2\nxx 3\nxx 4\nxx 5\nxx 6\nxx 7\nxx 8",
			},
		},

		{
			description: "2 elements, by size",
			rule: &config.Rule{
				Split: &config.Split{MaxSize: 20},
			},
			expect: []string{
				"xx 0\nxx 1\nxx 2\nxx 3", "xx 4\nxx 5\nxx 6\nxx 7", "xx 8",
			},
		},
		{
			description: "partition - filed index",
			rule: &config.Rule{
				Split: &config.Split{MaxSize: 1024,
					Partition: &config.Partition{
						FieldIndex: index,
						Hash:       "fnv",
						Mod:        2,
					},
				},
			},
			expect: []string{
				"xx 0\nxx 2\nxx 4\nxx 6\nxx 8", "xx 1\nxx 3\nxx 5\nxx 7",
			},
		},
	}

	for _, useCase := range useCases {
		var data = make([]string, 0)
		err := Split(strings.NewReader(text), func(partition interface{}) io.WriteCloser { return newTestWriter(&data) }, useCase.rule)
		assert.Nil(t, err)
		assert.EqualValues(t, useCase.expect, data, useCase.description)
	}

}

type testWriter struct {
	*bytes.Buffer
	data *[]string
}

func (t *testWriter) Close() error {
	*t.data = append(*t.data, t.String())
	return nil
}

func newTestWriter(data *[]string) io.WriteCloser {
	return &testWriter{
		data:   data,
		Buffer: new(bytes.Buffer),
	}
}
