package config

import (
	"fmt"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	"path"
	"smirror/base"
	"smirror/job"
	"strings"
	"time"
)

//Rule represent matching resource route rule
type Rule struct {
	Info    base.Info
	Dest    *Resource
	Source  *Resource
	Replace []*Replace `json:",omitempty"`

	Split *Split `json:",omitempty"`
	job.Actions
	*Compression
	//PreserveDepth  - preserves specified folder depth in dest URL
	PreserveDepth *int `json:",omitempty"`

	//Group defines group of rule to be matched, otherwise multi match is invalid
	Group string `json:",omitempty"`
}

//HasPreserveDepth returns true if property has been specified
func (r *Rule) HasPreserveDepth() bool {
	return r.PreserveDepth != nil
}

//GetPreserveDepth returns PreservceDepth
func (r *Rule) GetPreserveDepth() int {
	if r.PreserveDepth != nil {
		return *r.PreserveDepth
	}
	return 0
}

//Validate checks if route is valid
func (r *Rule) Validate() error {
	if r.Source == nil {
		return fmt.Errorf("source was empty")
	}
	if r.Dest == nil {
		return fmt.Errorf("dest was empty")
	}
	return nil
}

//HasMatch returns true if URL matches prefix or suffix
func (r *Rule) HasMatch(URL string) bool {
	if r.Source.Bucket != "" {
		bucket := url.Host(URL)
		if bucket != r.Source.Bucket {
			return false
		}
	}
	location := url.Path(URL)
	parent, name := path.Split(location)
	return r.Source.Match(parent, file.NewInfo(name, 0, 0644, time.Now(), false))
}

//Resources returns rule resource
func (r *Rule) Resources() []*Resource {
	var result = make([]*Resource, 0)
	if r.Source.Credentials != nil || r.Source.CustomKey != nil {
		result = append(result, r.Source)
	}
	if r.Dest.Credentials != nil || r.Dest.CustomKey != nil {
		result = append(result, r.Dest)
	}
	return result
}

//Name return route dest asset name
func (r *Rule) Name(URL string) string {
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

	depth := r.GetPreserveDepth()
	if depth == 0 && r.HasPreserveDepth() {
		return name
	}

	folderPath := strings.Trim(parent, "/")
	fragments := strings.Split(folderPath, "/")
	if !r.HasPreserveDepth() {
		depth = len(fragments)
	}

	fromRoot := false
	if depth < 0 {
		fromRoot = true
	}
	if depth < len(fragments) {
		if fromRoot {
			folderPath = strings.Join(fragments[depth*-1:], "/")
		} else {
			folderPath = strings.Join(fragments[len(fragments)-depth:], "/")
		}
	} else if strings.HasPrefix(folderPath, "/") {
		folderPath = string(folderPath[1:])
	}
	return path.Join(folderPath, name)
}
