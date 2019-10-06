package mon

import "smirror/config"

//Dataflow represents workflow info with unprocessed files
type Dataflow struct {
	config.Info
	ProcessedCount   int
	MaxProcessedSize int
	MinProcessedSize int
	UnprocessedCount int     `json:",omitempty"`
	DelaySec         int     `json:",omitempty"`
	Unprocessed      []*File `json:",omitempty"`
}

//NewWorkflow create a workflow
func NewWorkflow(info config.Info) *Dataflow {
	return &Dataflow{
		Info:        info,
		Unprocessed: make([]*File, 0),
	}
}
