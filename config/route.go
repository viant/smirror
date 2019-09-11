package config

import (
	"github.com/viant/afs/file"
	"github.com/viant/afs/matcher"
	"github.com/viant/afs/url"
	"path"
	"smirror/job"
	"strings"
	"time"
)

//Route represent matching resource route
type Route struct {
	Dest   Resource
	Source *Resource
	matcher.Basic
	Split *Split
	job.Actions
	*Compression
	//FolderDepth  - preserves specified folder depth in dest URL
	FolderDepth int
}

//HasMatch returns true if URL matches prefix or suffix
func (r *Route) HasMatch(URL string) bool {
	location := url.Path(URL)
	parent, name := path.Split(location)
	return r.Match(parent, file.NewInfo(name, 0, 0644, time.Now(), false))
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

	if r.FolderDepth == 0 {
		return name
	}

	folderPath := strings.Trim(parent, "/")
	fragments := strings.Split(folderPath, "/")
	if r.FolderDepth < len(fragments) {
		folderPath = strings.Join(fragments[len(fragments)-r.FolderDepth:], "/")
	} else if strings.HasPrefix(folderPath, "/") {
		folderPath = string(folderPath[1:])
	}
	return path.Join(folderPath, name)
}
