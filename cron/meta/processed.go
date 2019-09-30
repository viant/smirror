package meta

import (
	"time"
)

//Processed represents processed resource
type Processed struct {
	URL      string
	Modified time.Time
}

//NewProcessed create a processed resource state
func NewProcessed(URL string, modified time.Time) *Processed {
	return &Processed{URL: URL, Modified: modified}
}
