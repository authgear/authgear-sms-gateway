package twilio

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type TwilioClient struct {
	Client              *http.Client
	AccountSID          string
	AuthToken           string
	From                string
	MessagingServiceSID string
	Logger              *slog.Logger
}

func NewTwilioClient(httpClient *http.Client, accountSID string, authToken string, from string, messagingServiceSID string, logger *slog.Logger) *TwilioClient {
	return &TwilioClient{
		Client:              httpClient,
		AccountSID:          accountSID,
		AuthToken:           authToken,
		From:                from,
		MessagingServiceSID: messagingServiceSID,
		Logger:              logger,
	}
}

func (t *TwilioClient) send(options *smsclient.SendOptions) ([]byte, *SendResponse, error) {
	// Written against
	// https://www.twilio.com/docs/messaging/api/message-resource#create-a-message-resource

	u, err := url.Parse("https://api.twilio.com/2010-04-01/Accounts")
	if err != nil {
		return nil, nil, err
	}
	u = u.JoinPath(t.AccountSID, "Messages.json")

	values := url.Values{}
	values.Set("Body", options.Body)
	values.Set("To", string(options.To))

	if t.MessagingServiceSID != "" {
		values.Set("MessagingServiceSid", t.MessagingServiceSID)
	} else {
		values.Set("From", t.From)
	}

	requestBody := values.Encode()
	req, _ := http.NewRequest("POST", u.String(), strings.NewReader(requestBody))
	req.SetBasicAuth(t.AccountSID, t.AuthToken)

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	dumpedResponse, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, nil, err
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, errors.Join(
			err,
			&smsclient.SendResult{
				DumpedResponse: dumpedResponse,
			},
		)
	}

	sendResponse, err := ParseSendResponse(respData)
	if err != nil {
		return nil, nil, errors.Join(
			err,
			&smsclient.SendResult{
				DumpedResponse: dumpedResponse,
			},
		)
	}

	attrs := []slog.Attr{}
	if sendResponse.Status != nil {
		attrs = append(attrs, slog.String("status", *sendResponse.Status))
	}
	if sendResponse.SID != nil {
		attrs = append(attrs, slog.String("sid", *sendResponse.SID))
	}
	if sendResponse.DateCreated != nil {
		attrs = append(attrs, slog.String("date_created", *sendResponse.DateCreated))
	}
	if sendResponse.DateSent != nil {
		attrs = append(attrs, slog.String("date_sent", *sendResponse.DateSent))
	}
	if sendResponse.DateUpdated != nil {
		attrs = append(attrs, slog.String("date_updated", *sendResponse.DateUpdated))
	}
	if sendResponse.ErrorCode != nil {
		attrs = append(attrs, slog.Int("error_code", *sendResponse.ErrorCode))
	}
	if sendResponse.ErrorMessage != nil {
		attrs = append(attrs, slog.String("error_message", *sendResponse.ErrorMessage))
	}

	t.Logger.LogAttrs(context.TODO(), slog.LevelInfo, "twilio response", attrs...)

	return dumpedResponse, sendResponse, nil
}

func (t *TwilioClient) Send(options *smsclient.SendOptions) (*smsclient.SendResult, error) {
	info := &smsclient.SendResultInfo{
		SendResultInfoTwilio: &smsclient.SendResultInfoTwilio{},
	}
	info.SendResultInfoTwilio.BodyLength = len(options.Body)

	dumpedResponse, sendSMSResponse, err := t.send(options)
	if err != nil {
		return nil, err
	}

	var segmentCount *int
	if sendSMSResponse.NumSegments != nil {
		if parsed, err := strconv.Atoi(*sendSMSResponse.NumSegments); err == nil {
			segmentCount = &parsed
		}
	}
	info.SendResultInfoTwilio.SegmentCount = segmentCount

	return &smsclient.SendResult{
		DumpedResponse: dumpedResponse,
		Success:        sendSMSResponse.ErrorCode == nil,
		Info:           info,
	}, nil
}

var _ smsclient.RawClient = &TwilioClient{}
