package gs

import "fmt"

// Event is the payload of a GCS event.
type Event struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

//URL returns event source URL
func (e Event) URL() string {
	return fmt.Sprintf("gs://%v/%v", e.Bucket, e.Name)
}
