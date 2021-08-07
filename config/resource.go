package config

import (
	"fmt"
	"github.com/viant/afs/file"
	"github.com/viant/afs/matcher"
	"github.com/viant/afs/option"
	"github.com/viant/afs/url"
	"github.com/viant/toolbox/data"
	"github.com/viant/toolbox/data/udf"
	"regexp"
	"smirror/auth"
	"smirror/config/pattern"
	"strings"
)

//Represents a destination
type Resource struct {
	matcher.Basic
	Overflow    *Overflow
	Bucket      string            `json:",omitempty"`
	URL         string            `json:",omitempty"`
	Region      string            `json:",omitempty"`
	CustomKey   *CustomKey        `json:",omitempty"`
	Grant       *option.Grant     `json:",omitempty"`
	ACL         *option.ACL       `json:",omitempty"`
	Credentials *auth.Credentials `json:",omitempty"`
	Proxy       *option.Proxy
	Topic       string `json:",omitempty"`
	Queue       string `json:",omitempty"`
	Vendor      string `json:",omitempty"`
	//Optional pubsub project ID, otherwise it uses default one.
	ProjectID  string `json:",omitempty"`
	Pattern    string `json:",omitempty"`
	compiled   *regexp.Regexp
	Parameters []*pattern.Param `json:",omitempty"`
}

func (r *Resource) ExpandURL(sourceURL string) (string, error) {
	var err error
	if r.Pattern != "" && len(r.Parameters) > 0 {
		if r.compiled == nil {
			r.compiled, err = regexp.Compile(r.Pattern)
			if err != nil {
				return "", err
			}
		}
		var params = make(map[string]interface{})
		udfs := data.NewMap()
		udf.Register(udfs)
		for _, param := range r.Parameters {
			paramValue := expandWithPattern(r.compiled, sourceURL, param.Expression)
			params[param.Name] = udfs.ExpandAsText(paramValue)
		}
		expander := data.Map(params)
		return expander.ExpandAsText(r.URL), nil
	}
	return r.URL, nil
}

func expandWithPattern(expr *regexp.Regexp, sourceURL string, expression string) string {
	_, URLPath := url.Base(sourceURL, file.Scheme)
	matched := expr.FindStringSubmatch(URLPath)
	for i := 1; i < len(matched); i++ {
		key := fmt.Sprintf("$%v", i)
		count := strings.Count(expression, key)
		if count > 0 {
			expression = strings.Replace(expression, key, matched[i], count)
		}
	}
	return expression
}

//CloneWithURL clone resource with URL
func (r Resource) CloneWithURL(URL string) *Resource {
	return &Resource{
		Basic:       r.Basic,
		URL:         URL,
		Region:      r.Region,
		CustomKey:   r.CustomKey,
		Proxy:       r.Proxy,
		Grant:       r.Grant,
		Credentials: r.Credentials,
		Topic:       r.Topic,
		Queue:       r.Queue,
		ProjectID:   r.ProjectID,
	}
}

func (r *Resource) Init(projectID string) {
	if r.Topic == "" {
		return
	}
	if r.ProjectID == "" {
		r.ProjectID = projectID
	}
	if r.Topic != "" {
		if !strings.Contains(r.Topic, "/") && r.ProjectID != "" {
			r.Topic = fmt.Sprintf("projects/%s/topics/%s", r.ProjectID, r.Topic)
		}
	}
	if r.ProjectID == "" {
		if elements := strings.Split(r.Topic, ""); len(elements) == 4 {
			r.ProjectID = elements[1]
		}
	}
}
