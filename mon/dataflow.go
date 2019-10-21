package mon

import (
	"smirror/base"
)

//RuleInfo represents workflow info with unprocessed files
type RuleInfo struct {
	base.Info
	ProcessedCount   int
	MaxProcessedSize int
	MinProcessedSize int
	UnprocessedCount int     `json:",omitempty"`
	DelaySec         int     `json:",omitempty"`
	Unprocessed      []*File `json:",omitempty"`
}

//NewWorkflow create a workflow
func NewWorkflow(info base.Info) *RuleInfo {
	return &RuleInfo{
		Info:        info,
		Unprocessed: make([]*File, 0),
	}
}
