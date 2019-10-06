package config

import (
	"fmt"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	"path"
	"smirror/job"
	"strings"
	"time"
)

//Route represent matching resource route rule
type Route struct {
	Info    Info
	Dest    *Resource
	Source  *Resource
	Replace []*Replace `json:",omitempty"`

	Split *Split `json:",omitempty"`
	job.Actions
	*Compression
	//PreserveDepth  - preserves specified folder depth in dest URL
	PreserveDepth int `json:",omitempty"`
}

//Validate checks if route is valid
func (r *Route) Validate() error {
	if r.Source == nil {
		return fmt.Errorf("source was empty")
	}
	if r.Dest == nil {
		return fmt.Errorf("dest was empty")
	}
	return nil
}

//HasMatch returns true if URL matches prefix or suffix
func (r *Route) HasMatch(URL string) bool {
	location := url.Path(URL)
	parent, name := path.Split(location)
	return r.Source.Match(parent, file.NewInfo(name, 0, 0644, time.Now(), false))
}

//Name return route dest asset name
func (r *Route) Name(URL string) string {
	_, location := url.Base(URL, file.Scheme)
	parent, name := path.Split(location)
	ext := path.Ext(name)
	if r.Compression != nil {
		switch r.Compression.Codec {
		case GZipCodec:
			if ext != GZIPExtension {
				name += GZIPExtension
			}
		}
	} else {
		if ext == GZIPExtension {
			name = string(name[:len(name)-len(ext)])
		}
	}

	if r.PreserveDepth == 0 {
		return name
	}

	folderPath := strings.Trim(parent, "/")
	fragments := strings.Split(folderPath, "/")
	if r.PreserveDepth < len(fragments) {
		folderPath = strings.Join(fragments[len(fragments)-r.PreserveDepth:], "/")
	} else if strings.HasPrefix(folderPath, "/") {
		folderPath = string(folderPath[1:])
	}
	return path.Join(folderPath, name)
}
