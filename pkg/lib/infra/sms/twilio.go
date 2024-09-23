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

func (t *TwilioClient) Send(options *SendOptions) (*SendResult, error) {
	if t.TwilioClient == nil {
		return nil, ErrMissingTwilioConfiguration
	}

	params := &api.CreateMessageParams{}
	params.SetBody(options.Body)
	params.SetTo(options.To)
	if t.MessagingServiceSID != "" {
		params.SetMessagingServiceSid(t.MessagingServiceSID)
	} else {
		params.SetFrom(t.Sender)
	}

	resp, err := t.TwilioClient.Api.CreateMessage(params)
	if err != nil {
		return nil, fmt.Errorf("twilio: %w", err)
	}

	j, err := json.Marshal(resp)
	return &SendResult{
		ClientResponse: j,
	}, err
}

var _ RawClient = &TwilioClient{}
