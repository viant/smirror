package smirror

import (
	"io"
	"smirror/config"
	"strings"
	"unsafe"
)

//replacer represent replacer
type Rewriter struct {
	replacer *strings.Replacer
}

//HasReplacer returns true when replacer has replacement rules
func (t *Rewriter) HasReplacer() bool {
	return t.replacer != nil
}

func byteToString(data []byte) string {
	ptr := unsafe.Pointer(&data)
	return *(*string)(ptr)
}

func (t *Rewriter) Write(writer io.Writer, data []byte) error {
	if t.replacer == nil {
		_, err := writer.Write(data)
		return err
	}
	text := byteToString(data)
	_, err := t.replacer.WriteString(writer, text)
	return err
}

//New create a Rewriter
func NewRewriter(replaces ...*config.Replace) *Rewriter {
	if len(replaces) == 0 {
		return &Rewriter{}
	}
	pairs := make([]string, 0)
	for _, replace := range replaces {
		pairs = append(pairs, replace.From)
		pairs = append(pairs, replace.To)
	}
	return &Rewriter{replacer: strings.NewReplacer(pairs...)}
}
