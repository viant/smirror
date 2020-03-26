package option

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/viant/afs/url"
	"os"
	"path"
	"smirror"
	"smirror/base"
	"strings"
)

type Options struct {
	RuleURL string `short:"r" long:"rule" description:"rule URL"`

	Validate bool `short:"V" long:"validate" description:"run validation"`

	Version bool `short:"v" long:"version" description:"bqtail version"`

	SourceURL string `short:"s" long:"src" description:"source data URL" `

	DestinationURL string `short:"d" long:"dest" description:"destination table" `

	MatchPrefix string `short:"P" long:"prefix" description:"source match prefix"`

	MatchSuffix string `short:"S" long:"suffix" description:"source match suffix"`

	MatchPattern string `short:"R" long:"reg expr pattern" description:"source match reg expr pattern"`

	HistoryURL string `short:"H" long:"history" description:"history url to track already process file in previous run"`

	Stream bool `short:"X" long:"stream" description:"run constantly to stream changed/new datafile(s)"`

	ProjectID string `short:"p" long:"project" description:"Google Cloud Project"`

	PreserveDepth int `short:"D" long:"depth" description:"path preservation depth"`

	Topic string `short:"t" long:"topic" description:"Google Cloud Pub/Sub topic"`

	Queue string `short:"q" long:"queue" description:"AWS SQS queue"`
}

//HistoryPathURL return history URL
func (r *Options) HistoryPathURL(URL string) string {
	urlPath := url.Path(URL)
	historyName := md5Hash(urlPath) + ".json"
	historyName = strings.Replace(historyName, "=", "", strings.Count(historyName, "="))
	return url.Join(r.HistoryURL, historyName)
}

func (r *Options) Init(config *smirror.Config) {
	if r.SourceURL != "" {
		r.SourceURL = normalizeLocation(r.SourceURL)
	}
	if r.DestinationURL != "" {
		r.DestinationURL = normalizeLocation(r.DestinationURL)
	}
	if r.RuleURL != "" {
		r.RuleURL = normalizeLocation(r.RuleURL)
	}

	if r.HistoryURL != "" {
		r.HistoryURL = normalizeLocation(r.HistoryURL)
	}
	r.initHistoryURL()
}

//Hash returns fnv fnvHash value
func md5Hash(key string) string {
	h := md5.New()
	_, _ = h.Write([]byte(key))
	data := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

func (r *Options) initHistoryURL() {
	if r.HistoryURL != "" {
		return
	}
	r.HistoryURL = path.Join(os.Getenv("HOME"), ".smirror")
	if !r.Stream {
		r.HistoryURL = url.Join(base.InMemoryStorageBaseURL, r.HistoryURL)
	}
}

func normalizeLocation(location string) string {
	if location == "" {
		return ""
	}
	if strings.HasPrefix(location, "~/") {
		location = strings.Replace(location, "~/", os.Getenv("HOME"), 1)
	}

	if url.Scheme(location, "") == "" && !strings.HasPrefix(location, "/") {
		currentDirectory, _ := os.Getwd()
		fmt.Printf("%v\n", currentDirectory)
		return path.Join(currentDirectory, location)
	}
	fmt.Printf("normailized %v\n", location)
	return location
}
