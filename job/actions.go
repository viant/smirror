package job

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/smirror/base"
)

//Actions represents a job completion
type Actions struct {
	OnSuccess []*Action
	OnFailure []*Action
}

//Run run completion
func (a *Actions) Run(context *Context, service afs.Service, notify Notify, info *base.Info, body interface{}) error {
	err := a.run(context, service, notify, info, body)
	if err != nil {
		JSON, _ := json.Marshal(a)
		err = errors.Wrapf(err, "failed to run post actions: %s", JSON)
	}
	return err
}

//Run run completion
func (a *Actions) run(context *Context, service afs.Service, notify Notify, info *base.Info, body interface{}) error {
	actions := a.OnSuccess
	isError := context.Error != nil
	if context.Error != nil {
		actions = a.OnFailure
	}
	if len(actions) == 0 {
		return nil
	}
	for _, action := range actions {
		e := action.Do(context, service, notify, info, body)
		if e == nil && isError {
			e = action.WriteError(context, service)
		}
		if e != nil {
			return e
		}
	}
	return nil
}
