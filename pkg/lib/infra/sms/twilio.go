package sms

import (
	"errors"
	"fmt"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

var ErrMissingTwilioConfiguration = errors.New("twilio: configuration is missing")

type TwilioClient struct {
	Name                string
	TwilioClient        *twilio.RestClient
	Sender              string
	MessagingServiceSID string
}

func NewTwilioClient(name string, accountSID string, authToken string, sender string, messagingServiceSID string) *TwilioClient {
	return &TwilioClient{
		Name: name,
		TwilioClient: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: accountSID,
			Password: authToken,
		}),
		Sender:              sender,
		MessagingServiceSID: messagingServiceSID,
	}
}

func (t *TwilioClient) GetName() string {
	return t.Name
}

func (t *TwilioClient) Send(
	to string,
	body string,
	templateName string,
	languageTag string,
	templateVariables *TemplateVariables,
) error {
	if t.TwilioClient == nil {
		return ErrMissingTwilioConfiguration
	}

	params := &api.CreateMessageParams{}
	params.SetBody(body)
	params.SetTo(to)
	if t.MessagingServiceSID != "" {
		params.SetMessagingServiceSid(t.MessagingServiceSID)
	} else {
		params.SetFrom(t.Sender)
	}

	_, err := t.TwilioClient.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("twilio: %w", err)
	}

	return nil
}

var _ RawClient = &TwilioClient{}
