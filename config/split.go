package config

import (
	"fmt"
	"path"
	"strings"
)

//Split represents a split rule
type Split struct {
	//MaxLines max number lines in one split chunk
	MaxLines int
	//Template has to have %v placeholder for file name, and %d (or padded placeholder i.e. %04d) chunk number
	Template string

	//MaxSize max size, if file larger then splits
	MaxSize int
}

//Name returns a chunk name for supplied URL and mirrorChunkeddAsset number
func (s *Split) Name(router *Route, URL string, counter int32) string {
	name := router.Name(URL)
	destName := ""
	ext := ""
	if extIndex := strings.Index(name, "."); extIndex != -1 {
		ext = string(name[extIndex+1:])
		name = string(name[:extIndex])
	}

	parent, child := path.Split(name)
	if s.Template != "" {
		lastIndex := strings.LastIndex(s.Template, "%")
		nameIndex := strings.Index(s.Template, "%v")
		if nameIndex == lastIndex {
			destName = fmt.Sprintf(s.Template, counter, child)
		} else {
			destName = fmt.Sprintf(s.Template, child, counter)
		}
	} else {
		destName = fmt.Sprintf("%04d_%v", counter, child)
	}
	destName += "." + ext

	return path.Join(parent, destName)
}
