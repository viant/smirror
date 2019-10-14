package mon

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/afs/storage"
	"github.com/viant/toolbox"
	"smirror/base"
	"smirror/config"
	"strings"
	"time"
)

const defaultTriggerAge = "1hour"
const defaultErrorRecency = "1hour"
const agoKeyword = "Ago"

//Request represents monitoring request
type Request struct {
	UnprocessedDuration       string
	unprocessedModifiedBefore *time.Time

	ProcessedURL           string
	ProcessedRecency       string
	processedModifiedAfter *time.Time

	ErrorRecency       string
	errorModifiedAfter *time.Time

	ConfigURL  string
	TriggerURL string
	ErrorURL   string
}

//Response represents monitoring response
type Response struct {
	UnprocessedCount int
	MaxDelayInSec    int `json:",omitempty"`
	ProcessCount     int `json:",omitempty"`
	ProcessedBytes   int `json:",omitempty"`
	ErrorCount       int
	Dataflows        []*Dataflow `json:",omitempty"`
	Errors           []*Error    `json:",omitempty"`
	Status           string
	Error            string `json:",omitempty"`
	workflowMap      map[string]*Dataflow
}

func (r *Response) AddError(object storage.Object, message string) {
	mirrorError := &Error{URL: object.URL(), Message: message, Created: object.ModTime()}
	r.Errors = append(r.Errors, mirrorError)
	r.ErrorCount++
}

func (r *Response) AddUnprocessed(now time.Time, route *config.Rule, file storage.Object) {
	info := base.Info{
		Workflow: base.UnclassifiedDataflow,
	}
	if route != nil {
		info = route.Info
	}
	workflow, ok := r.workflowMap[info.Workflow]
	if !ok {
		workflow = NewWorkflow(info)
		r.Dataflows = append(r.Dataflows, workflow)
		r.workflowMap[info.Workflow] = workflow
	}
	workflow.UnprocessedCount++
	r.UnprocessedCount++
	elapsed := now.Sub(file.ModTime())
	workflow.Unprocessed = append(workflow.Unprocessed, &File{
		URL:      file.URL(),
		Modified: file.ModTime(),
		Size:     int(file.Size()),
		Age:      fmt.Sprintf("%s", ((elapsed / time.Second) * time.Second)),
	})
}

func (r *Response) AddProcessed(route *config.Rule, object storage.Object) {
	info := base.Info{
		Workflow: base.UnclassifiedDataflow,
	}
	if route != nil {
		info = route.Info
	}
	workflow, ok := r.workflowMap[info.Workflow]
	if !ok {
		workflow = NewWorkflow(info)
		r.Dataflows = append(r.Dataflows, workflow)
		r.workflowMap[info.Workflow] = workflow
	}
	workflow.ProcessedCount++
	r.ProcessCount++
	fileSize := int(object.Size())
	if fileSize > workflow.MaxProcessedSize {
		workflow.MaxProcessedSize = fileSize
	}
	if fileSize < workflow.MinProcessedSize || workflow.MinProcessedSize == 0 {
		workflow.MinProcessedSize = fileSize
	}
	r.ProcessedBytes += fileSize
}

//Init initialises request
func (r *Request) Init() (err error) {
	if r.UnprocessedDuration == "" {
		r.UnprocessedDuration = defaultTriggerAge
	}
	if !(strings.Contains(strings.ToLower(r.UnprocessedDuration), "ago") || strings.Contains(strings.ToLower(r.UnprocessedDuration), "past")) {
		r.UnprocessedDuration += agoKeyword
	}
	if r.unprocessedModifiedBefore, err = toolbox.TimeAt(r.UnprocessedDuration); err != nil {
		return errors.Wrapf(err, "invalid UnprocessedDuration: %v", r.UnprocessedDuration)
	}

	if r.ProcessedRecency == "" {
		r.ProcessedRecency = defaultTriggerAge
	}
	if !(strings.Contains(strings.ToLower(r.ProcessedRecency), "ago") || strings.Contains(strings.ToLower(r.ProcessedRecency), "past")) {
		r.ProcessedRecency += agoKeyword
	}
	if r.processedModifiedAfter, err = toolbox.TimeAt(r.ProcessedRecency); err != nil {
		return errors.Wrapf(err, "invalid ErrorRecency: %v", r.ProcessedRecency)
	}

	if r.ErrorRecency == "" {
		r.ErrorRecency = defaultErrorRecency
	}
	if !(strings.Contains(strings.ToLower(r.ErrorRecency), "ago") || strings.Contains(strings.ToLower(r.ErrorRecency), "past")) {
		r.ErrorRecency += agoKeyword
	}
	if r.errorModifiedAfter, err = toolbox.TimeAt(r.ErrorRecency); err != nil {
		return errors.Wrapf(err, "invalid ErrorRecency: %v", r.ErrorRecency)
	}
	return err
}

//Validate check if request is valid
func (r *Request) Validate() (err error) {
	if r.ConfigURL == "" {
		return errors.Errorf("configURL was empty")
	}
	if r.TriggerURL == "" {
		return errors.Errorf("triggerURL was empty")
	}
	return nil
}

//NewResponse create a response
func NewResponse() *Response {
	return &Response{
		Status:      base.StatusOK,
		workflowMap: make(map[string]*Dataflow),
		Dataflows:   make([]*Dataflow, 0),
		Errors:      make([]*Error, 0),
	}
}
