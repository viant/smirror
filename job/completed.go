package job

type Completion struct {
	OnSuccess []*Action
	OnFailure []*Action
}

func (c *Completion) Run(context *Context) error {
	actions := c.OnSuccess
	if context.Error != nil {
		actions = c.OnFailure
	}
	if len(actions) == 0 {
		return nil
	}
	for _, action := range actions {
		if err := action.Do(context);err != nil {
			return err
		}
	}
	return nil
}