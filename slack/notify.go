package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"smirror/job"
)

func (s *service) Notify(ctx context.Context, request *job.NotifyRequest) error {
	err := s.notify(ctx, request)
	if err != nil {
		err = errors.Wrapf(err, "failed to notify on slack: %v", request.Channels)
		fmt.Printf("%v\n", err)
	}
	return err
}

func (s *service) notify(ctx context.Context, request *job.NotifyRequest) error {
	err := request.Init(s.Region, s.projectID)
	if err == nil {
		if request.Credentials == nil {
			request.Credentials = s.Credentials
		}
		err = request.Validate()
	}
	if err != nil {
		return err
	}

	if request.Credentials.RawToken == "" {
		data, err := s.Secret.Decrypt(ctx, &request.Credentials.Secret)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(data, &request.OAuthToken); err != nil {
			return errors.Wrapf(err, "failed to unmarshal token: %s", data)
		}
		request.RawToken  = request.Token
		request.Token = ""
	}
	client := slack.New(request.RawToken)
	return s.postMessage(ctx, client, request)
}

func (s *service) uploadFile(context context.Context, client *slack.Client, request *job.NotifyRequest) error {
	if request.Body == nil {
		return nil
	}
	body := ""
	switch value := request.Body.(type) {
	case []byte:
		body = string(value)
	case string:
		body = value
	default:

		data, err := json.MarshalIndent(request.Body, "", "\t")
		if err != nil {
			err = errors.Wrapf(err, "failed to decode body: %v", value)
			return err
		}
		body = string(data)
	}
	if body == "" {
		return nil
	}

	fileType := "text"
	if json.Valid([]byte(body)) {
		fileType = "json"
	}
	uploadRequest := slack.FileUploadParameters{
		Filename: request.Filename,
		Title:    request.Title,
		Filetype: fileType,
		Content:  string(body),
		Channels: request.Channels,
	}
	_, err := client.UploadFile(uploadRequest)
	return err
}

func (s *service) postMessage(context context.Context, client *slack.Client, request *job.NotifyRequest) error {
	err := s.sendMessage(context, client, request)
	if err == nil {
		err = s.uploadFile(context, client, request)
	}
	return err
}

func (s *service) sendMessage(context context.Context, client *slack.Client, request *job.NotifyRequest) (err error) {
	if request.Message == "" {
		return nil
	}
	attachment := slack.Attachment{
		Text:       request.Message,
		AuthorName: request.From,
	}
	for _, channel := range request.Channels {
		if _, _, e := client.PostMessage(channel, slack.MsgOptionText(request.Title, false), slack.MsgOptionAttachments(attachment)); e != nil {
			err = e
		}
	}
	return err
}
