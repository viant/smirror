package mirror

//Datafile represents a data file
type Datafile struct {
	URL string
	*Compression
}

//NewDatafile returns a new datafile
func NewDatafile(URL string, compression *Compression) *Datafile {
	return &Datafile{URL: URL, Compression: compression}
}
