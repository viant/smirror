package cmd

import (
	"context"
	"github.com/pkg/errors"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	"smirror/cmd/validate"
	"smirror/shared"
)



func (s *service) Validate(ctx context.Context, request *validate.Request) error {
	request.Init(s.config)
	if request.RuleURL == "" {
		return errors.Errorf("ruleURL was empty")
	}
	parent, _ := url.Split(request.RuleURL, file.Scheme)
	cfg, err := newConfig(ctx, s.config.ProjectID)
	if err != nil {
		return errors.Wrap(err, "failed to create config for validation")
	}
	cfg.Mirrors.BaseURL = parent
	err = cfg.Init(ctx, s.fs)
	if err == nil && len(cfg.Mirrors.Rules) > 0 {
		s.reportRule(cfg.Mirrors.Rules[0])
		shared.LogLn("Rule is VALID\n")
	}

	return err
}
