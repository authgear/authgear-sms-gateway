package sms

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

var ErrMissingTwilioConfiguration = errors.New("twilio: configuration is missing")

type TwilioClient struct {
	TwilioClient        *twilio.RestClient
	Sender              string
	MessagingServiceSID string
}

func NewTwilioClient(accountSID string, authToken string, sender string, messagingServiceSID string) *TwilioClient {
	return &TwilioClient{
		TwilioClient: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: accountSID,
			Password: authToken,
		}),
		Sender:              sender,
		MessagingServiceSID: messagingServiceSID,
	}
}

func (t *TwilioClient) Send(
	to string,
	body string,
	templateName string,
	languageTag string,
	templateVariables *TemplateVariables,
) (ClientResponse, error) {
	if t.TwilioClient == nil {
		return []byte{}, ErrMissingTwilioConfiguration
	}

	params := &api.CreateMessageParams{}
	params.SetBody(body)
	params.SetTo(to)
	if t.MessagingServiceSID != "" {
		params.SetMessagingServiceSid(t.MessagingServiceSID)
	} else {
		params.SetFrom(t.Sender)
	}

	resp, err := t.TwilioClient.Api.CreateMessage(params)
	if err != nil {
		return ClientResponse{}, fmt.Errorf("twilio: %w", err)
	}

	j, err := json.Marshal(resp)
	return ClientResponse(j), err
}

var _ RawClient = &TwilioClient{}
