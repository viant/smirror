package smirror

import "github.com/viant/smirror/config"

//Datafile represents a data file
type Datafile struct {
	URL string
	*config.Compression
}

//CompressionCodec returns destination code
func (d *Datafile) CompressionCodec() string {
	if d.Compression == nil {
		return ""
	}
	return d.Compression.Codec
}

//NewDatafile returns a new datafile
func NewDatafile(URL string, compression *config.Compression) *Datafile {
	return &Datafile{URL: URL, Compression: compression}
}
