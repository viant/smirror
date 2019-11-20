package smirror

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/option"
	"github.com/viant/afs/storage"
	"github.com/viant/afs/url"
	"github.com/viant/afsc/gs"
	"github.com/viant/afsc/s3"
	"io"
	"io/ioutil"
	"os"
	"path"
	"smirror/base"
	"smirror/config"
	"smirror/contract"
	"smirror/job"
	"smirror/msgbus"
	"smirror/msgbus/pubsub"
	"smirror/msgbus/sqs"
	"smirror/secret"
	"smirror/slack"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

//Slack represents a mirror service
type Service interface {
	//Mirror copies/split source to matched destination
	Mirror(ctx context.Context, request *contract.Request) *contract.Response
}

type service struct {
	mux      *sync.Mutex
	config   *Config
	fs       afs.Service
	secret   secret.Service
	msgbus   msgbus.Service
	notifier slack.Slack
}

func (s *service) Mirror(ctx context.Context, request *contract.Request) *contract.Response {

	request.Attempt++
	response := contract.NewResponse(request.URL)

	err := s.mirror(ctx, request, response)
	if err != nil {
		response.Status = base.StatusError
		response.Error = err.Error()
	}
	if response.Error == "" {
		return response
	}

	if IsNotFound(response.Error) {
		response.Status = base.StatusNoFound
		response.Error = ""
		response.NotFoundError = response.Error
	} else if IsRetryError(response.Error) {
		if request.Attempt < s.config.MaxRetries {
			return s.Mirror(ctx, request)
		}
	}
	return response
}

func (s *service) mirror(ctx context.Context, request *contract.Request, response *contract.Response) (err error) {
	_, err = s.config.Mirrors.ReloadIfNeeded(ctx, s.fs)
	if err != nil {
		return err
	}
	var rule *config.Rule
	matched := s.config.Mirrors.Match(request.URL)
	switch len(matched) {
	case 0:
	case 1:
		rule = matched[0]
	default:
		JSON, _ := json.Marshal(matched)
		return errors.Errorf("multi rule match currently not supported: %s", JSON)
	}
	response.TotalRules = len(s.config.Mirrors.Rules)
	if rule == nil {
		response.Status = base.StatusNoMatch
		return nil
	}
	if err := s.initRule(ctx, rule); err != nil {
		return errors.Wrapf(err, "railed to initialise rule: %v", rule.Info.Workflow)
	}
	response.Rule = rule
	options, err := s.secret.StorageOpts(ctx, rule.Source.CloneWithURL(request.URL))
	if err != nil {
		return err
	}
	object, err := s.fs.Object(ctx, request.URL, options...)
	if object == nil {
		response.Status = base.StatusNoFound
		response.NotFoundError = fmt.Sprintf("does not exist: %v", err)
		return nil
	}
	response.ChecksumSkip = int(object.Size()) > s.config.Streaming.ChecksumSkipThreshold
	if int(object.Size()) > s.config.Streaming.Threshold {
		response.StreamOption = option.NewStream(s.config.Streaming.PartSize, int(object.Size()))
	}

	err = s.mirrorAsset(ctx, rule, request.URL, response)
	jobContent := job.NewContext(ctx, err, request.URL, response.Rule.Name(request.URL))
	response.TimeTakenMs = int(time.Now().Sub(request.Timestamp) / time.Millisecond)
	if e := rule.Actions.Run(jobContent, s.fs, s.notifier.Notify, &response.Rule.Info, response); e != nil && err == nil {
		err = e
	}
	return err
}

func (s *service) addStreamingOptions(options []storage.Option, streamOpt *option.Stream) []storage.Option {
	if streamOpt != nil {
		options = append(options, streamOpt)
	}
	return options
}

func (s *service) mirrorAsset(ctx context.Context, rule *config.Rule, URL string, response *contract.Response) error {
	transferStream := s.transferStream
	if rule.Split != nil {
		transferStream = s.transferChunkStream
	}
	options, err := s.secret.StorageOpts(ctx, rule.Source.CloneWithURL(URL))
	if err != nil {
		return errors.Wrapf(err, "failed to get storage option for %v", rule.Source)
	}
	options = s.addStreamingOptions(options, response.StreamOption)
	if rule.ShallArchiveWalk(URL) {
		archvieURL := rule.ArchiveWalkURL(URL)
		return s.fs.Walk(ctx, archvieURL, func(ctx context.Context, baseURL string, parent string, info os.FileInfo, reader io.Reader) (toContinue bool, err error) {
			if info.IsDir() {
				return true, nil
			}
			location := path.Join(parent, info.Name())
			streamURL := url.Join(URL, location)
			err = transferStream(ctx, ioutil.NopCloser(reader), streamURL, rule, response)
			return err == nil, err
		}, options...)
	}
	reader, err := s.fs.DownloadWithURL(ctx, URL, options...)
	if err != nil {
		return errors.Wrapf(err, "failed to download source: %v", URL)
	}
	defer func() {
		_ = reader.Close()
	}()
	return transferStream(ctx, reader, URL, rule, response)
}

func (s *service) transferStream(ctx context.Context, reader io.ReadCloser, URL string, rule *config.Rule, response *contract.Response) (err error) {
	sourceCompression := rule.SourceCompression(URL)
	reader, err = NewReader(reader, sourceCompression)
	if err != nil {
		return errors.Wrapf(err, "failed to create reader")
	}
	destName := rule.Name(URL)
	destURL := url.Join(rule.Dest.URL, destName)
	destCompression := rule.Compression

	if path.Ext(URL) == path.Ext(destURL) && len(rule.Replace) == 0 {
		sourceCompression = nil
		destCompression = nil
	}
	dataCopy := &Transfer{
		skipChecksum: response.ChecksumSkip,
		rewriter:     NewRewriter(rule.Replace...),
		Resource:     rule.Dest,
		Reader:       reader,
		Dest:         NewDatafile(destURL, destCompression)}

	return s.transfer(ctx, dataCopy, response)
}

func (s *service) transferChunkStream(ctx context.Context, reader io.ReadCloser, URL string, rule *config.Rule, response *contract.Response) (err error) {
	sourceCompression := rule.SourceCompression(URL)
	reader, err = NewReader(reader, sourceCompression)
	if err != nil {
		return errors.Wrapf(err, "failed to create reader")
	}
	counter := int32(0)
	waitGroup := &sync.WaitGroup{}
	rewriter := NewRewriter(rule.Replace...)
	err = Split(reader, s.chunkWriter(ctx, URL, rule, &counter, waitGroup, response), rule.Split, rewriter)
	if err == nil {
		waitGroup.Wait()
	}
	return err
}

func (s *service) transfer(ctx context.Context, transfer *Transfer, response *contract.Response) error {
	if transfer.Resource.Topic != "" || transfer.Resource.Queue != "" {
		return s.publish(ctx, transfer, response)
	}
	if transfer.Resource.URL != "" {
		err := s.upload(ctx, transfer, response)
		if err != nil {
			return errors.Wrapf(err, "failed to transfer to: %v", transfer.Dest.URL)
		}
		return nil
	}
	JSON, _ := json.Marshal(transfer)
	return fmt.Errorf("invalid transfer: %s", JSON)
}

func (s *service) publish(ctx context.Context, transfer *Transfer, response *contract.Response) error {
	reader, err := transfer.GetReader()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	switch s.config.SourceScheme {
	case gs.Scheme, s3.Scheme:
		attributes := make(map[string]interface{})
		attributes[base.SourceAttribute] = transfer.Dest.URL
		dest := transfer.Resource.Topic
		if dest == "" {
			dest = transfer.Resource.Queue
		}
		dest = strings.Replace(dest, "$partition", transfer.partition, 1)
		pubResponse, err := s.msgbus.Publish(ctx, &msgbus.Request{
			Dest:       dest,
			Data:       data,
			Attributes: attributes,
		})
		if err != nil {
			if IsNotFound(err.Error()) {
				return errors.Errorf("failed to publish data, no such topic: %v", transfer.Resource.Topic)
			}
			return err
		}
		response.MessageIDs = append(response.MessageIDs, pubResponse.MessageIDs...)
		return nil
	}
	return fmt.Errorf("unsupported message msgbus %v", s.config.SourceScheme)
}

func (s *service) upload(ctx context.Context, transfer *Transfer, response *contract.Response) error {
	reader, err := transfer.GetReader()
	if err != nil {
		return errors.Wrapf(err, "failed to get reader for: %v", transfer.Resource.URL)
	}
	options, err := s.secret.StorageOpts(ctx, transfer.Resource)
	if err != nil {
		return err
	}
	if transfer.skipChecksum {
		options = append(options, option.NewSkipChecksum(true))
	}
	if err = s.fs.Upload(ctx, transfer.Dest.URL, file.DefaultFileOsMode, reader, options...); err != nil {
		return errors.Wrapf(err, "failed to upload %v", transfer.Dest.URL)
	}
	response.AddURL(transfer.Dest.URL)
	return nil
}

//Init initialises this service
func (s *service) Init(ctx context.Context) error {
	return s.config.Init(ctx, s.fs)
}

func (s *service) initActions(actions []*job.Action) {
	if len(actions) == 0 {
		return
	}
	for i := range actions {
		if actions[i].Action == job.ActionNotify && actions[i].Credentials == nil {
			actions[i].Credentials = s.config.SlackCredentials
		}
	}
}

//initRule updates resources
func (s *service) initRule(ctx context.Context, rule *config.Rule) (err error) {
	resources := rule.Resources()
	s.initActions(rule.OnSuccess)
	s.initActions(rule.OnFailure)
	if len(resources) > 0 {
		if err = s.secret.Init(ctx, s.fs, resources); err != nil {
			return errors.Wrap(err, "failed to init resource secrets")
		}
	}

	if s.config.UseMessageDest() && s.msgbus == nil {
		switch s.config.SourceScheme {
		case gs.Scheme:
			if s.msgbus, err = pubsub.New(ctx, s.config.ProjectID); err != nil {
				return errors.Wrapf(err, "unable to create pubsub publisher for %v", s.config.SourceScheme)
			}
		case s3.Scheme:
			if s.msgbus, err = sqs.New(ctx); err != nil {
				return errors.Wrapf(err, "unable to create sqs publisher for %v", s.config.SourceScheme)
			}
		default:
			return errors.Errorf("unsupported scheme for publisher: '%v'", s.config.SourceScheme)
		}
	}
	return nil
}

func (s *service) chunkWriter(ctx context.Context, URL string, rule *config.Rule, counter *int32, waitGroup *sync.WaitGroup, response *contract.Response) func(partition interface{}) io.WriteCloser {
	return func(partition interface{}) io.WriteCloser {
		splitCount := atomic.AddInt32(counter, 1)
		destName := rule.Split.Name(rule, URL, splitCount, partition)
		destURL := url.Join(rule.Dest.URL, destName)
		return NewWriter(rule, func(writer *Writer) error {
			if writer.Reader == nil {
				return fmt.Errorf("Writer reader was empty")
			}
			waitGroup.Add(1)
			defer waitGroup.Done()
			dataCopy := &Transfer{
				partition:    fmt.Sprintf("%v", partition),
				skipChecksum: response.ChecksumSkip,
				Resource:     rule.Dest,
				Reader:       writer.Reader,
				Dest:         NewDatafile(destURL, nil),
			}
			return s.transfer(ctx, dataCopy, response)
		})
	}
}

//NewSlack creates a new mirror service
func New(ctx context.Context, config *Config) (Service, error) {
	err := config.Init(ctx, afs.New())
	if err != nil {
		return nil, err
	}
	fs := afs.New()
	secretService := secret.New(config.SourceScheme, fs)
	result := &service{config: config,
		fs:       fs,
		mux:      &sync.Mutex{},
		secret:   secretService,
		notifier: slack.NewSlack(config.Region, config.ProjectID, fs, secretService, config.SlackCredentials),
	}
	return result, result.Init(ctx)
}
