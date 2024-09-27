package nexmo

import (
	"encoding/json"
	"errors"
	"fmt"

	nexmo "github.com/njern/gonexmo"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
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

func (n *NexmoClient) Send(options *smsclient.SendOptions) (*smsclient.SendResult, error) {
	if n.NexmoClient == nil {
		return nil, ErrMissingNexmoConfiguration
	}

	message := nexmo.SMSMessage{
		From:  n.Sender,
		To:    string(options.To),
		Type:  nexmo.Text,
		Text:  options.Body,
		Class: nexmo.Standard,
	}

	resp, err := n.NexmoClient.SMS.Send(&message)
	if err != nil {
		return nil, fmt.Errorf("nexmo: %w", err)
	}

	if resp.MessageCount == 0 {
		err = errors.New("nexmo: no sms is sent")
		return nil, err
	}

	report := resp.Messages[0]
	if report.ErrorText != "" {
		err = fmt.Errorf("nexmo: %s", report.ErrorText)
		return nil, err
	}

	j, err := json.Marshal(resp)
	return &smsclient.SendResult{
		ClientResponse: j,
		Success:        report.Status == nexmo.ResponseSuccess,
	}, err
}

var _ smsclient.RawClient = &NexmoClient{}