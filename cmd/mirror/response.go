package mirror

import (
	"smirror/base"
	"sync"
	"sync/atomic"
)

//Response represents response
type Response struct {
	Status   string
	Errors   []string
	DataURLs []string
	Mirrored []string
	NoMatch []string
	Failed []string
	historyURLs []string
	mux sync.Mutex
	pending int32
}

//History returns events history
func (r Response) HistoryURLs() []string {
	return r.historyURLs
}

//AddHistoryURL adds history URL
func (r *Response) AddHistoryURL(URL string) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.historyURLs = append(r.historyURLs, URL)
}

//IncrementPending increments pending
func (r *Response) IncrementPending(deleta int32) {
	atomic.AddInt32(&r.pending, deleta)
}

//Pending returns pending events
func (r Response) Pending() int {
	return int(atomic.LoadInt32(&r.pending))
}

//AddMirrored adds trigger URL
func (r *Response) AddMirrored(URL string) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.Mirrored = append(r.Mirrored, URL)
}


//AddFailed adds trigger URL
func (r *Response) AddFailed(URL string) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.Failed = append(r.Failed, URL)
}

//AddNoMatch adds trigger URL
func (r *Response) AddNoMatch(URL string) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.NoMatch = append(r.NoMatch, URL)
}


//AddDataURLs adds trigger URL
func (r *Response) AddDataURLs(URL string) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.DataURLs = append(r.DataURLs, URL)
}



//NewResponse creates a response
func NewResponse() *Response {
	return &Response{
		Status: base.StatusOK,
		Errors: make([]string, 0),
		DataURLs:make([]string, 0),
		Mirrored:make([]string, 0),
		NoMatch:make([]string, 0),
		Failed:make([]string, 0),
	}
}


//AddError adds response error
func (r *Response) AddError(err error) {
	if err == nil {
		return
	}
	r.mux.Lock()
	defer r.mux.Unlock()
	r.Status = base.StatusError
	r.Errors = append(r.Errors, err.Error())
}