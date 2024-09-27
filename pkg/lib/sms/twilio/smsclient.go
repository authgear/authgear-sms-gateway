package twilio

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
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

func (t *TwilioClient) Send(options *smsclient.SendOptions) (*smsclient.SendResult, error) {
	if t.TwilioClient == nil {
		return nil, ErrMissingTwilioConfiguration
	}

	params := &api.CreateMessageParams{}
	params.SetBody(options.Body)
	params.SetTo(string(options.To))
	if t.MessagingServiceSID != "" {
		params.SetMessagingServiceSid(t.MessagingServiceSID)
	} else {
		params.SetFrom(t.Sender)
	}

	resp, err := t.TwilioClient.Api.CreateMessage(params)
	if err != nil {
		return nil, fmt.Errorf("twilio: %w", err)
	}

	var segmentCount *int
	if resp.NumSegments != nil {
		if i, err := strconv.Atoi(*resp.NumSegments); err == nil {
			segmentCount = &i
		}
	}

	j, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return &smsclient.SendResult{
		// FIXME: Switch to call Twilio via REST.
		DumpedResponse: j,
		Success:        resp.ErrorCode != nil,
		SegmentCount:   segmentCount,
	}, err
}

var _ smsclient.RawClient = &TwilioClient{}
