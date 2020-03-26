package cmd

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/viant/afs/file"
	"github.com/viant/afs/mem"
	"github.com/viant/afs/option"
	"github.com/viant/afs/url"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"smirror/base"
	"smirror/cmd/build"
	"smirror/config"
)

func (s *service) Build(ctx context.Context, request *build.Request) error {
	hasMatcher := request.MatchPrefix != "" || request.MatchSuffix != "" || request.MatchPattern != ""
	hasRuleURL := request.RuleURL != ""
	request.Init(s.config)
	if ! hasRuleURL {
		request.RuleURL = url.Join(ruleBaseURL, "rule.yaml")
	}

	rule := &config.Rule{
		Dest: &config.Resource{
			URL: request.DestinationURL,
		},
		Source: &config.Resource{
		},
		Info: base.Info{
			URL:          request.RuleURL,
			LeadEngineer: os.Getenv("USER"),
		},
	}

	if request.MatchPrefix != "" {
		rule.Source.Prefix = request.MatchPrefix
	}
	if request.MatchSuffix != "" {
		rule.Source.Suffix = request.MatchSuffix
	}
	if request.MatchPattern != "" {
		rule.Source.Filter = request.MatchPattern
	}

	if request.SourceURL != "" && ! hasMatcher {
		if files, _ := s.fs.List(ctx, request.SourceURL, option.NewRecursive(true)); len(files) > 0 {
			URLPath := url.Path(files[0].URL())
			rule.Source.Prefix, _ = path.Split(URLPath)
			rule.Source.Suffix = path.Ext(files[0].Name())
		}
	}
	rule.Streaming = &config.Streaming{
		ThresholdMb:             300,
		PartSizeMb:              15,
		ChecksumSkipThresholdMb: 400,
	}

	if request.Topic != "" {
		rule.Dest.Topic = request.Topic
	}
	if request.Queue != "" {
		rule.Dest.Queue = request.Queue
	}

	rule.PreserveDepth = &request.PreserveDepth
	if hasRuleURL {
		s.reportRule(rule)
		return nil
	}

	ruleMap := ruleToMap(rule)
	ruleYAML, err := yaml.Marshal(ruleMap)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal rule")
	}

	if mem.Scheme == url.Scheme(request.RuleURL, "") {
		err = s.fs.Upload(ctx, rule.Info.URL, file.DefaultFileOsMode, bytes.NewReader(ruleYAML))
	}
	return err
}
