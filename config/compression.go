package config

import "strings"

const (
	//GZipCodec gzip code
	GZipCodec = "gzip"
	//GZIPExtension gzip extension
	GZIPExtension = ".gz"

	//ZipCodec zip code
	ZipCodec = "zip"
	//ZIPExtension zip extension
	ZIPExtension = ".zip"

	//TarCodec tar code
	TarCodec = "tar"
	//TarExtension tar extension
	TarExtension = ".tar"
)

//Compression represents conversion strategy
type Compression struct {
	Codec      string `json:",omitempty"`
	Uncompress bool   `json:",omitempty"`
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
