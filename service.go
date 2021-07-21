package smirror

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/cache"
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
	"smirror/shared"
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
	mux          *sync.Mutex
	config       *Config
	fs           afs.Service
	cfs          afs.Service
	secret       secret.Service
	msgbus       msgbus.Service
	msgbusVendor string
	notifier     slack.Slack
}

func (s *service) Mirror(ctx context.Context, request *contract.Request) *contract.Response {

	request.Attempt++
	response := contract.NewResponse(request.URL)

	err := s.mirror(ctx, request, response)
	if err != nil {
		response.Status = base.StatusError
		response.Error = err.Error()
	}
	if s.config.ResponseURL != "" {
		s.logResponse(ctx, response)
	}
	if response.Error == "" {
		return response
	}
	if IsNotFound(response.Error) {
		response.Status = base.StatusNoFound
		response.NotFoundError = response.Error
		response.Error = ""
	} else if IsRetryError(response.Error) {
		if request.Attempt < s.config.MaxRetries {
			return s.Mirror(ctx, request)
		}
	}
	if s.config.ResponseURL != "" {
		s.logResponse(ctx, response)
	}
	return response
}

func (s *service) mirror(ctx context.Context, request *contract.Request, response *contract.Response) (err error) {
	_, err = s.config.Mirrors.ReloadIfNeeded(ctx, s.cfs)
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
	if rule.Disabled {
		response.Status = base.StatusDisabled
		return nil
	}

	if err := s.initRule(ctx, rule); err != nil {
		return errors.Wrapf(err, "frailed to initialise rule: %v", rule.Info.Workflow)
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
	response.FileSize = object.Size()
	if rule.Source.Overflow != nil {
		if rule.Source.Overflow.Size() < object.Size() {
			return s.handleOverflow(ctx, object, rule.Source.Overflow, rule, request, response)
		}
	}

	if rule.DoneMarker != "" {
		parentURL, _ := url.Split(request.URL, file.Scheme)
		if object.Name() == rule.DoneMarker {
			return s.replay(ctx, parentURL, rule.DoneMarker)
		}
		//Check if marker file is present, otherwise delay transfer
		markerURL := url.Join(parentURL, rule.DoneMarker)
		if exists, _ := s.fs.Exists(ctx, markerURL, option.NewObjectKind(true)); !exists {
			response.Status = base.StatusPartial
			return nil
		}
	}

	var streaming = &s.config.Streaming
	if rule.Streaming != nil {
		streaming = rule.Streaming
	}
	response.ChecksumSkip = int(object.Size()) > streaming.ChecksumSkipThreshold()
	if int(object.Size()) > streaming.Threshold() {
		response.StreamOption = option.NewStream(streaming.PartSize(), int(object.Size()))
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
	reader, err := s.fs.OpenURL(ctx, URL, options...)
	if err != nil {
		return errors.Wrapf(err, "failed to download source: %v", URL)
	}
	defer func() {
		_ = reader.Close()
	}()
	return transferStream(ctx, reader, URL, rule, response)
}

func (s *service) transferStream(ctx context.Context, reader io.Reader, URL string, rule *config.Rule, response *contract.Response) (err error) {
	reader, err = NewReader(rule, reader, response, URL)
	if err != nil {
		return errors.Wrapf(err, "failed to create reader")
	}
	destName := rule.Name(URL)
	baseDestURL, err := rule.Dest.ExpandURL(URL)
	if err != nil {
		return errors.Wrapf(err, "failed to expanded URL")
	}
	destURL := url.Join(baseDestURL, destName)
	destCompression := rule.Compression
	if path.Ext(URL) == path.Ext(destURL) && !rule.HasTransformer() {
		destCompression = nil
	}

	dataCopy := &Transfer{
		skipChecksum: response.ChecksumSkip,
		stream:       rule.Streaming,
		rule:         rule,
		Resource:     rule.Dest,
		Reader:       reader,
		Dest:         NewDatafile(destURL, destCompression)}

	return s.transfer(ctx, dataCopy, response)
}

func (s *service) transferChunkStream(ctx context.Context, reader io.Reader, URL string, rule *config.Rule, response *contract.Response) (err error) {
	reader, err = NewReader(rule, reader, response, URL)
	if err != nil {
		return errors.Wrapf(err, "failed to create reader")
	}
	counter := int32(0)
	waitGroup := &sync.WaitGroup{}
	err = Split(reader, s.chunkWriter(ctx, URL, rule, &counter, waitGroup, response), rule)
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
		if base.IsSchemaError(err) {
			response.SchemaError = err.Error()
		}
		if err != nil {
			return errors.Wrapf(err, "failed to transfer to: %v", transfer.Dest.URL)
		}
		return nil
	}
	JSON, _ := json.Marshal(transfer)
	return fmt.Errorf("dest.URL was empty: invalid transfer: %s: ", JSON)
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

	switch s.msgbusVendor {
	case shared.VendorPubsub, shared.VendorSQS:
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
	return fmt.Errorf("unsupported message vendor %v", s.msgbusVendor)
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
		if stream := transfer.stream; stream != nil && stream.PartSizeMb > 0 {
			options = append(options, option.NewStream(stream.PartSize(), int(response.FileSize)))
		}
	}
	if rule := transfer.rule; rule != nil && rule.AllowEmpty {
		options = append(options, option.NewEmpty(rule.AllowEmpty))
	}
	writer, err := s.fs.NewWriter(ctx, transfer.Dest.URL, file.DefaultFileOsMode, options...)
	if err != nil {
		return err
	}
	response.AddURL(transfer.Dest.URL)
	if transfer.Dest.CompressionCodec() == config.GZipCodec {
		gzipWriter := gzip.NewWriter(writer)
		if _, err = io.Copy(gzipWriter, reader); err != nil {
			return err
		}
		if err = gzipWriter.Flush(); err == nil {
			err = gzipWriter.Close()
		}

	} else {
		if _, err = io.Copy(writer, reader); err != nil {
			return err
		}
	}
	err = writer.Close()
	if err != nil {
		//if errors mirroring delete dest corrupted transfer
		s.fs.Delete(ctx, transfer.Dest.URL)
	}
	return err
}

//Load initialises this service
func (s *service) Init(ctx context.Context) error {
	return s.config.Init(ctx, s.cfs)
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
		if rule.Dest.Vendor == "" {
			switch s.config.SourceScheme {
			case gs.Scheme:
				rule.Dest.Vendor = shared.VendorPubsub
			case s3.Scheme:
				rule.Dest.Vendor = shared.VendorSQS
			default:
				if rule.Dest.Topic != "" {
					rule.Dest.Vendor = shared.VendorPubsub
				}
				if rule.Dest.Queue != "" {
					rule.Dest.Vendor = shared.VendorSQS
				}
			}

		}
		s.msgbusVendor = rule.Dest.Vendor
		switch rule.Dest.Vendor {
		case shared.VendorPubsub:
			if s.config.ProjectID == "" {
				s.config.ProjectID = os.Getenv("GCLOUD_PROJECT")
			}

			if s.msgbus, err = pubsub.New(ctx, s.config.ProjectID); err != nil {
				return errors.Wrapf(err, "unable to create pubsub publisher for %v", rule.Dest.Vendor)
			}
		case shared.VendorSQS:
			if s.msgbus, err = sqs.New(ctx); err != nil {
				return errors.Wrapf(err, "unable to create sqs publisher for %v", rule.Dest.Vendor)
			}
		default:
			return errors.Errorf("unsupported message bus vendor: '%v'", rule.Dest.Vendor)
		}
	}
	return nil
}

func (s *service) chunkWriter(ctx context.Context, URL string, rule *config.Rule, counter *int32, waitGroup *sync.WaitGroup, response *contract.Response) func(partition interface{}) io.WriteCloser {
	return func(partition interface{}) io.WriteCloser {
		splitCounter := atomic.AddInt32(counter, 1)
		destName := rule.Split.Name(rule, URL, splitCounter, partition)
		return NewWriter(rule, func(writer *Writer) error {
			baseDestURL, err := rule.Dest.ExpandURL(URL)
			if err != nil {
				return fmt.Errorf("failed to expand URL due to %w", err)
			}
			destURL := url.Join(baseDestURL, destName)
			if writer.Reader == nil {
				return fmt.Errorf("reader was empty")
			}
			waitGroup.Add(1)
			defer waitGroup.Done()
			dataCopy := &Transfer{
				rule:         rule,
				splitCounter: splitCounter,
				partition:    fmt.Sprintf("%v", partition),
				skipChecksum: response.ChecksumSkip,
				stream:       rule.Streaming,
				Resource:     rule.Dest,
				Reader:       writer.Reader,
				Dest:         NewDatafile(destURL, nil),
			}
			return s.transfer(ctx, dataCopy, response)
		})
	}
}

func (s *service) replay(ctx context.Context, parentURL, doneMarker string) error {
	objects, err := s.fs.List(ctx, parentURL)
	if err != nil {
		return err
	}
	replayer := base.NewReplayer(s.fs)
	replayer.Run(ctx, 5)
	for _, object := range objects {
		if object.IsDir() || object.Name() == doneMarker {
			continue
		}
		replayer.Schedule(object.URL())
	}
	return replayer.Wait()
}

func (s *service) logResponse(ctx context.Context, response *contract.Response) {
	if response.Rule != nil {
		response.RuleURL = response.Rule.Info.URL
	}

	//avoid event infinitive cycle
	if strings.Contains(response.TriggeredBy, s.config.ResponseURL) {
		return
	}
	response.Rule = nil
	JSON, err := json.Marshal(response)
	if err != nil {
		response.LogError = err.Error()
		return
	}
	UUID := uuid.New().String()
	logURL := url.Join(s.config.ResponseURL, UUID+".json")
	err = s.fs.Upload(ctx, logURL, file.DefaultFileOsMode, bytes.NewReader(JSON))
	if err != nil {
		response.LogError = err.Error()
	}
}

func (s *service) handleOverflow(ctx context.Context, object storage.Object, overflow *config.Overflow, rule *config.Rule, request *contract.Request, response *contract.Response) error {
	response.Status = base.StatusOverflow
	_, URLPath := url.Base(object.URL(), file.Scheme)
	destURL := url.Join(overflow.DestURL, URLPath)
	err := s.fs.Copy(ctx, object.URL(), destURL) //change to move
	if err != nil {
		response.Error = err.Error()
		return err
	}
	response.DestURLs = append(response.DestURLs, destURL)
	msgService, err := s.overflowBusService(ctx, overflow)
	if err != nil || msgService == nil {
		return fmt.Errorf("msg bus service was empty")
	}
	data, err := json.Marshal(overflow.MessageEvent(destURL))
	if err != nil {
		return fmt.Errorf("failed to marshal message: %+v", overflow.MessageEvent(destURL))
	}
	msg := &msgbus.Request{
		Dest: overflow.MessageDest(),
		Data: data,
	}
	output, err := msgService.Publish(ctx, msg)
	if err != nil {
		return err
	}
	response.MessageIDs = output.MessageIDs
	jobContent := job.NewContext(ctx, err, request.URL, response.Rule.Name(request.URL))
	response.TimeTakenMs = int(time.Now().Sub(request.Timestamp) / time.Millisecond)
	if e := rule.Actions.Run(jobContent, s.fs, s.notifier.Notify, &response.Rule.Info, response); e != nil && err == nil {
		err = e
	}
	return err
}

func (s *service) overflowBusService(ctx context.Context, overflow *config.Overflow) (msgbus.Service, error) {
	var service msgbus.Service
	var err error
	if overflow.Queue != "" {
		service, err = sqs.New(ctx)
	} else if overflow.Topic != "" {
		service, err = pubsub.New(ctx, overflow.ProjectID)
	}
	return service, err
}

//NewSlack creates a new mirror service
func New(ctx context.Context, config *Config) (Service, error) {
	cfs := cache.Singleton(config.URL)
	err := config.Init(ctx, cfs)
	if err != nil {
		return nil, err
	}
	fs := afs.New()
	secretService := secret.New(config.SourceScheme, fs)
	result := &service{config: config,
		fs:       fs,
		cfs:      cfs,
		mux:      &sync.Mutex{},
		secret:   secretService,
		notifier: slack.NewSlack(config.Region, config.ProjectID, fs, secretService, config.SlackCredentials),
	}
	return result, result.Init(ctx)
}
