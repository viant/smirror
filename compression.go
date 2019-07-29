package smirror

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

//NewCompressionForURL returns compression for matched codec or nil
func NewCompressionForURL(URL string) *Compression {
	if strings.HasSuffix(URL, GZIPExtension) {
		return &Compression{Codec: GZipCodec}
	}
	return nil
}
