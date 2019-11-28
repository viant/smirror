package config

import (
	"compress/gzip"
	"fmt"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	"io"
	"path"
	"smirror/base"
	"smirror/job"
	"strings"
	"time"
)

//Rule represent matching resource route rule
type Rule struct {
	Info    base.Info
	Disabled bool  `json:",omitempty"`
	Dest    *Resource
	Source  *Resource
	Replace []*Replace `json:",omitempty"`
	Recover *Recover

	Streaming *Streaming
	Split *Split `json:",omitempty"`
	job.Actions
	*Compression
	//PreserveDepth  - preserves specified folder depth in dest URL
	PreserveDepth *int `json:",omitempty"`

	//Group defines group of rule to be matched, otherwise multi match is invalid
	Group string `json:",omitempty"`
}

//NewReplacer create a replaced for the rule
func (r *Rule) NewReplacer() *strings.Replacer {
	if len(r.Replace) == 0 {
		return nil
	}
	pairs := make([]string, 0)
	for _, replace := range r.Replace {
		pairs = append(pairs, replace.From)
		pairs = append(pairs, replace.To)
	}
	return strings.NewReplacer(pairs...)
}



//NewReader returns a reader for a rule
func (r *Rule) NewReader(reader io.Reader, sourceURL string) (io.Reader, error) {
	compression := r.SourceCompression(sourceURL)
	var err error
	if compression != nil && compression.Codec == GZipCodec {
		if reader, err =  gzip.NewReader(reader);err != nil {
			return reader, err
		}
	}
	if ! r.HasTransformer() {
		 return reader, nil
	}
	return NewTransformer(reader, r)
}


//HasTransformer returns true if rule has recover or replace option
func (r *Rule) HasTransformer() bool {
	return r.Recover != nil || len(r.Replace) > 0
}

//HasSplit returns true if rule has split defined
func (r *Rule) HasSplit() bool {
	if r.Split == nil {
		return false
	}
	return r.Split.MaxSize > 0 || r.Split.MaxLines > 0 || r.Split.Partition != nil
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

func (r *Rule) ShallArchiveWalk(URL string) bool {
	if r.Compression == nil {
		return false
	}
	return (strings.HasSuffix(URL, TarExtension) || strings.HasSuffix(URL, ZIPExtension)) && r.Compression.Uncompress
}

func (r *Rule) ArchiveWalkURL(URL string) string {
	ext := path.Ext(URL)
	ext = strings.Replace(ext, ".", "", 1)
	return fmt.Sprintf("%v/%v://localhost/", strings.Replace(URL, "://", ":", 1), ext)
}

//Validate checks if route is valid
func (r *Rule) Validate() error {
	r.Init()
	if r.Source == nil {
		return fmt.Errorf("source was empty")
	}
	if r.Dest == nil {
		return fmt.Errorf("dest was empty")
	}
	return nil
}

//Init initialises routes
func (r *Rule) Init() {
	if r.HasSplit() || r.HasTransformer() {
		if r.Compression == nil {
			r.Compression = &Compression{}
		}
		r.Compression.Uncompress = true
	}
	if r.Streaming != nil {
		r.Streaming.Init()
	}
}

//SourceCompression returns compression for URL
func (r *Rule) SourceCompression(URL string) (source *Compression) {
	source = NewCompressionForURL(URL)
	compression := r.Compression
	if compression != nil && compression.Uncompress {
		if source != nil {
			source.Uncompress = compression.Uncompress
		}
		return source
	}
	hasDestCompression := (compression != nil && compression.Codec != "") && source != nil
	if (hasDestCompression && source.Equals(compression)) || !hasDestCompression {
		return nil
	}
	return source
}

//Match returns true if URL matches prefix or suffix
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
	sourceCompression := r.SourceCompression(URL)
	_, location := url.Base(URL, file.Scheme)
	parent, name := path.Split(location)
	ext := path.Ext(name)
	if r.Compression != nil && r.Compression.Codec != "" {
		switch r.Compression.Codec {
		case GZipCodec:
			if ext != GZIPExtension {
				name += GZIPExtension
			}
		}
	} else if sourceCompression != nil && sourceCompression.Uncompress {
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
		depth = -1 * depth
		fromRoot = true
	}
	if depth <= len(fragments) {
		if fromRoot {
			folderPath = strings.Join(fragments[depth:], "/")
		} else {
			folderPath = strings.Join(fragments[len(fragments)-depth:], "/")
		}
	} else if strings.HasPrefix(folderPath, "/") {
		folderPath = string(folderPath[1:])
	}
	return path.Join(folderPath, name)
}
