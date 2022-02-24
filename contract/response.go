package contract

import (
	"github.com/viant/afs/option"
	"github.com/viant/smirror/base"
	"github.com/viant/smirror/config"

	"sync"
	"time"
)

//Response represents a response
type Response struct {
	TriggeredBy   string
	FileSize      int64 `json:",omitempty"`
	LogError      string `json:",omitempty"`
	DestURLs      []string `json:",omitempty"`
	MessageIDs    []string `json:",omitempty"`
	TimeTakenMs   int
	Rule          *config.Rule `json:",omitempty"`
	RuleURL       string
	TotalRules    int
	Status        string
	Error         string `json:",omitempty"`
	SchemaError   string `json:",omitempty"`
	NotFoundError string `json:",omitempty"`
	StartTime     time.Time
	BadRecords    int            `json:",omitempty"`
	ChecksumSkip  bool           `json:",omitempty"`
	StreamOption  *option.Stream `json:",omitempty"`
	mutex         *sync.Mutex
}

//AddURL adds url to dest urls
func (r *Response) AddURL(URL string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.DestURLs = append(r.DestURLs, URL)
}

//NewResponse returns a new response
func NewResponse(triggeredBy string) *Response {
	return &Response{
		Status:      base.StatusOK,
		TriggeredBy: triggeredBy,
		StartTime:   time.Now(),
		DestURLs:    make([]string, 0),
		MessageIDs:  make([]string, 0),
		mutex:       &sync.Mutex{},
	}
}
