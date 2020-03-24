package cmd

import (
	"bytes"
	"context"
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
	request.Init(s.config)
	if request.RuleURL == "" {
		request.RuleURL = url.Join(ruleBaseURL, "rule.yaml")
	}
	rule := &config.Rule{
		Dest: &config.Resource{
			URL: request.DestinationURL,
		},
		Source: &config.Resource{
			URL: request.SourceURL,
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

	hasMatcher := rule.Source.Prefix != "" || rule.Source.Suffix != "" || rule.Source.Filter != ""
	if request.SourceURL != "" && hasMatcher {
		if files, _ := s.fs.List(ctx, request.SourceURL, option.NewRecursive(true)); len(files) > 0 {
			rule.Source.Prefix, _ = url.Split(files[0].URL(), file.Scheme)
			rule.Source.Suffix = path.Ext(files[0].Name())
		}
	}
	rule.Streaming = &config.Streaming{
		ThresholdMb:             300,
		PartSizeMb:              15,
		ChecksumSkipThresholdMb: 400,
	}

	if !(request.SourceURL != "" || request.Validate) {
		s.reportRule(rule)
		return nil
	}
	ruleMap := ruleToMap(rule)
	ruleYAML, err := yaml.Marshal(ruleMap)
	if err != nil {
		return err
	}
	if mem.Scheme == url.Scheme(rule.Info.URL, "") {
		err = s.fs.Upload(ctx, rule.Info.URL, file.DefaultFileOsMode, bytes.NewReader(ruleYAML))
	}
	return err
}
