package job

import "github.com/viant/afs"

//Completion represents a job completion
type Completion struct {
	OnSuccess []*Action
	OnFailure []*Action
}

//Run run completion
func (c *Completion) Run(context *Context, service afs.Service) error {
	actions := c.OnSuccess
	isError := context.Error != nil
	if context.Error != nil {
		actions = c.OnFailure
	}
	if len(actions) == 0 {
		return nil
	}
	for _, action := range actions {
		err := action.Do(context, service)
		if err == nil && isError {
			err = action.WriteError(context, service)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
