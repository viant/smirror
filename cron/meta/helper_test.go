package meta

import (
	"github.com/viant/afs/file"
	"github.com/viant/afs/mem"
	"github.com/viant/afs/object"
	"github.com/viant/afs/storage"
	"github.com/viant/afs/url"
	"time"
)

func GetTestObjects(baseURL string, objects map[string]time.Time) []storage.Object {
	var result = make([]storage.Object, 0)
	for name, mod := range objects {
		_, name := url.Split(name, mem.Scheme)
		info := file.NewInfo(name, 0, 0644, mod, false)
		result = append(result, object.New(url.Join(baseURL, name), info, nil))
	}
	return result
}
