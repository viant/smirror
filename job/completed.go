package job

//Completion represents a job completion
type Completion struct {
	OnSuccess []*Action
	OnFailure []*Action
}

//Run run completion
func (c *Completion) Run(context *Context) error {
	actions := c.OnSuccess
	isError := context.Error != nil
	if context.Error != nil {
		actions = c.OnFailure
	}
	if len(actions) == 0 {
		return nil
	}
	for _, action := range actions {
		err := action.Do(context)
		if err == nil && isError {
			err = action.WriteError(context)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
