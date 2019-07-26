package mirror

import (
	"fmt"
	"strings"
)

//Split represents a split rule
type Split struct {
	//MaxLines max number lines in one split chunk
	MaxLines int
	//Template has to have %v placeholder for file name, and %d (or padded placeholder i.e. %04d) chunk number
	Template string
}

//Name returns a chunk name for supplied URL and mirrorSplittedAsset number
func (s *Split) Name(router *Route, URL string, counter int32) string {
	name := router.Name(URL)
	destName := ""
	ext := ""
	if extIndex := strings.Index(name, "."); extIndex != -1 {
		ext = string(name[extIndex+1:])
		name = string(name[:extIndex])
	}
	if s.Template != "" {
		lastIndex := strings.LastIndex(s.Template, "%")
		nameIndex := strings.Index(s.Template, "%v")
		if nameIndex == lastIndex {
			destName = fmt.Sprintf(s.Template, counter, name)
		} else {
			destName = fmt.Sprintf(s.Template, name, counter)
		}
	} else {
		destName = fmt.Sprintf("%04d_%v", counter, name)
	}
	destName += "." + ext
	return destName
}
