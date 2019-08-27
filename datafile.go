package smirror

import "smirror/config"

//Datafile represents a data file
type Datafile struct {
	URL string
	*config.Compression
}

//NewDatafile returns a new datafile
func NewDatafile(URL string, compression *config.Compression) *Datafile {
	return &Datafile{URL: URL, Compression: compression}
}
