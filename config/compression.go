package config

import "strings"

const (
	//GZipCodec gzip code
	GZipCodec = "gzip"
	//GZIPExtension gzip extension
	GZIPExtension = ".gz"
)

//Compression represents conversion strategy
type Compression struct {
	Codec string
}

//Equals returns true if compression is the same
func (c *Compression) Equals(compression *Compression) bool {
	if c == nil {
		if compression == nil {
			return true
		}
		return false
	}
	if compression == nil {
		return false
	}
	return c.Codec == compression.Codec
}

//NewCompressionForURL returns compression for matched codec or nil
func NewCompressionForURL(URL string) *Compression {
	if strings.HasSuffix(URL, GZIPExtension) {
		return &Compression{Codec: GZipCodec}
	}
	return nil
}
