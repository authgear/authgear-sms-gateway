package sms

import (
	"errors"
	"fmt"

	nexmo "github.com/njern/gonexmo"
)

var ErrMissingNexmoConfiguration = errors.New("nexmo: configuration is missing")

type NexmoClient struct {
	Name        string
	NexmoClient *nexmo.Client
	Sender      string
}

func NewNexmoClient(name string, apiKey string, apiSecret string, sender string) *NexmoClient {
	nexmoClient, _ := nexmo.NewClient(apiKey, apiSecret)
	return &NexmoClient{
		Name:        name,
		NexmoClient: nexmoClient,
		Sender:      sender,
	}
}

func (n *NexmoClient) GetName() string {
	return n.Name
}

func (n *NexmoClient) Send(
	to string,
	body string,
	templateName string,
	languageTag string,
	templateVariables *TemplateVariables,
) error {
	if n.NexmoClient == nil {
		return ErrMissingNexmoConfiguration
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
		return fmt.Errorf("nexmo: %w", err)
	}

	if resp.MessageCount == 0 {
		err = errors.New("nexmo: no sms is sent")
		return err
	}

	report := resp.Messages[0]
	if report.ErrorText != "" {
		err = fmt.Errorf("nexmo: %s", report.ErrorText)
		return err
	}

	return nil
}

var _ RawClient = &NexmoClient{}
