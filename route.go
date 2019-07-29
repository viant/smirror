package mirror

import (
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/url"
	"path"
	"smirror/job"
	"strings"
)

//Route represent matching resource route
type Route struct {
	Prefix       string
	Suffix       string
	DestURL      string
	Split        *Split
	OnCompletion job.Completion
	*Compression

	//FolderDepth  - preserves specified folder depth in dest URL
	FolderDepth int
}

//Routes represents route slice
type Routes []*Route

//HasMatch returns true if URL matches prefix or suffix
func (r *Route) HasMatch(URL string) bool {
	resource := url.NewResource(URL)
	location := resource.ParsedURL.Path
	if r.Prefix != "" {
		if !strings.HasPrefix(location, r.Prefix) {
			return false
		}
	}
	if r.Suffix != "" {
		if !strings.HasSuffix(location, r.Suffix) {
			return false
		}
	}
	return true
}

//Name return route dest asset name
func (r *Route) Name(URL string) string {
	parent, name := toolbox.URLSplit(URL)

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
	resource := url.NewResource(parent)
	folderPath := resource.ParsedURL.Path
	fragments := strings.Split(folderPath, "/")
	if r.FolderDepth < len(fragments) {
		folderPath = strings.Join(fragments[len(fragments)-r.FolderDepth:], "/")
	} else if strings.HasPrefix(folderPath, "/") {
		folderPath = string(folderPath[1:])
	}
	return path.Join(folderPath, name)
}

//HasMatch returns the first match route
func (r Routes) HasMatch(URL string) *Route {
	for i := range r {
		if r[i].HasMatch(URL) {
			return r[i]
		}
	}
	return nil
}
