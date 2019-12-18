package smirror

import (
	"compress/gzip"
	"io"
	"smirror/config"
	"smirror/config/recover"
)

//NewReader returns a reader for a rule
func NewReader(rule *config.Rule, reader io.Reader, sourceURL string) (io.Reader, error) {
	compression := rule.SourceCompression(sourceURL)
	var err error
	if compression != nil && compression.Codec == config.GZipCodec {
		if reader, err = gzip.NewReader(reader); err != nil {
			return reader, err
		}
	}
	if !rule.HasTransformer() {
		return reader, nil
	}
	if rule.Recover != nil || len(rule.Replace) > 0 {
		if reader, err = recover.NewReader(reader, rule); err != nil {
			return nil, err
		}
	}
	return reader, err
}
