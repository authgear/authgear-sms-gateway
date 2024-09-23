package sms

import (
	"encoding/json"
	"errors"
	"fmt"

	nexmo "github.com/njern/gonexmo"
)

var ErrMissingNexmoConfiguration = errors.New("nexmo: configuration is missing")

type NexmoClient struct {
	NexmoClient *nexmo.Client
	Sender      string
}

func NewNexmoClient(apiKey string, apiSecret string, sender string) *NexmoClient {
	nexmoClient, _ := nexmo.NewClient(apiKey, apiSecret)
	return &NexmoClient{
		NexmoClient: nexmoClient,
		Sender:      sender,
	}
}

func (n *NexmoClient) Send(
	to string,
	body string,
	templateName string,
	languageTag string,
	templateVariables *TemplateVariables,
) (ClientResponse, error) {
	if n.NexmoClient == nil {
		return ClientResponse{}, ErrMissingNexmoConfiguration
	}

	message := nexmo.SMSMessage{
		From:  n.Sender,
		To:    to,
		Type:  nexmo.Text,
		Text:  body,
		Class: nexmo.Standard,
	}

	resp, err := n.NexmoClient.SMS.Send(&message)
	if err != nil {
		return ClientResponse{}, fmt.Errorf("nexmo: %w", err)
	}

	if resp.MessageCount == 0 {
		err = errors.New("nexmo: no sms is sent")
		return ClientResponse{}, err
	}

	report := resp.Messages[0]
	if report.ErrorText != "" {
		err = fmt.Errorf("nexmo: %s", report.ErrorText)
		return ClientResponse{}, err
	}

	j, err := json.Marshal(resp)
	return ClientResponse(j), err
}

var _ RawClient = &NexmoClient{}
