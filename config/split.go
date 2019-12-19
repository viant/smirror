package config

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/twmb/murmur3"
	"github.com/viant/toolbox"
	"hash"
	"hash/fnv"
	"path"
	"strings"
)

//Split represents a split rule
type Split struct {
	//MaxLines max number lines in one split chunk
	MaxLines int
	//Template has to have %s placeholder for file name, and %d (or padded placeholder i.e. %04d) chunk number, %v is for partition
	Template string

	//Partition partition rule
	Partition *Partition

	//Schema format specific schema
	Schema string

	//SchemaURL format specific schema location
	SchemaURL string

	//MaxSize max size, if file larger then splits
	MaxSize int
}

//Partition represent partition split
type Partition struct {
	Field       string
	FieldIndex  int
	Separator   string
	Mod         int
	Hash        string
	keyProvider func(data []byte) (interface{}, error)
}

func newKeyProvider(partition *Partition) func(data []byte) (interface{}, error) {

	if partition.Field != "" {
		aMap := map[string]interface{}{}
		field := partition.Field
		return func(data []byte) (interface{}, error) {
			aMap[field] = ""
			err := json.Unmarshal(data, &aMap)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to unmarshal JSON")
			}
			key, _ := aMap[field]
			return key, nil
		}
	}

	separator := partition.Separator
	if separator == "" {
		separator = ","
	}
	index := partition.FieldIndex
	return func(data []byte) (interface{}, error) {
		elements := strings.Split(string(data), separator)
		if index >= len(elements) {
			return "", errors.Errorf("index out of bound, index: %v, len: %v, %s", index, len(elements), data)
		}
		value := elements[index]
		if len(value) > 0 {
			value = strings.Trim(value, `'"\n\t `)
		}
		return value, nil
	}
}

func bytesToInt(data []byte) int {
	keyNumeric := int64(0)
	shift := 0
	for i := 0; i < 8 && i < len(data); i++ {
		v := int64(data[len(data)-1-i])
		if shift == 0 {
			keyNumeric |= v
		} else {
			keyNumeric |= v << uint64(shift)
		}
		shift += 8
	}
	if keyNumeric < 0 {
		keyNumeric *= -1
	}
	return int(keyNumeric)
}

func (p *Partition) Key(data []byte) (interface{}, error) {
	if p.keyProvider == nil {
		p.keyProvider = newKeyProvider(p)
	}
	key, err := p.keyProvider(data)
	if err != nil {
		return nil, err
	}

	if p.Hash == "" && p.Mod == 0 {
		return key, nil
	}
	keyNumeric := toolbox.AsInt(key)
	if p.Hash != "" {
		var keyHash hash.Hash
		switch strings.ToLower(p.Hash) {
		case "murmur":
			keyHash = murmur3.New64()
		case "fnv":
			keyHash = fnv.New64()
		case "md5":
			keyHash = md5.New()
		default:
			return "", errors.Errorf("unsupported hash: %v", p.Hash)
		}

		keyData := toolbox.AsString(key)
		if _, err = keyHash.Write([]byte(keyData)); err != nil {
			return "", err
		}
		keyNumeric = bytesToInt(keyHash.Sum(nil))
	}

	if p.Mod == 0 {
		return keyNumeric, nil
	}
	return keyNumeric % p.Mod, nil
}

//Name returns a chunk name for supplied URL and mirrorChunkeddAsset number
func (s *Split) Name(router *Rule, URL string, counter int32, partition interface{}) string {
	name := router.Name(URL)
	destName := ""
	ext := ""
	if extIndex := strings.Index(name, "."); extIndex != -1 {
		ext = string(name[extIndex+1:])
		name = string(name[:extIndex])
	}
	parent, child := path.Split(name)

	templ := s.Template
	if templ == "" {
		templ = "%04d_%s"
		if s.Partition != nil {
			templ += "_%v"
		}
	}

	if templ != "" {
		templ = strings.Replace(templ, "%s", "$name", 1)
		templ = strings.Replace(templ, "%v", "$partition", 1)
		templ = strings.Replace(templ, "$chunk", "%03d", 1)
	}

	destName = templ
	destName = strings.Replace(destName, "$name", child, 2)
	if partition != nil && strings.Contains(destName, "$partition") {
		destName = strings.Replace(destName, "$partition", toolbox.AsString(partition), 1)
	}
	if strings.Contains(destName, "%") {
		destName = fmt.Sprintf(destName, counter)
	}
	destName += "." + ext
	return path.Join(parent, destName)
}
