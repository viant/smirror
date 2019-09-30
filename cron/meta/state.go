package meta

import (
	"github.com/viant/afs/storage"
	"time"
)

//State meta files storing processed resources
type State struct {
	Processed []*Processed
}

//Add adds storage object to processed resources
func (s *State) Add(objects ...storage.Object) {
	if len(s.Processed) == 0 {
		s.Processed = make([]*Processed, 0)
	}
	for i := range objects {
		s.Processed = append(s.Processed, NewProcessed(objects[i].URL(), objects[i].ModTime()))
	}
}

//Prune removes any resourced older than supplied max age
func (s *State) Prune(now time.Time, maxAge time.Duration) {
	if maxAge <= 0 {
		return
	}
	var survivors = make([]*Processed, 0)
	for i := range s.Processed {
		age := now.Sub(s.Processed[i].Modified)
		if age > maxAge {
			continue
		}
		survivors = append(survivors, s.Processed[i])
	}
	s.Processed = survivors
}

//ProcessMap returns processed resource map
func (s *State) ProcessMap() map[string]time.Time {
	var result = make(map[string]time.Time)
	if len(s.Processed) == 0 {
		return result
	}
	for i := range s.Processed {
		result[s.Processed[i].URL] = s.Processed[i].Modified
	}
	return result
}
