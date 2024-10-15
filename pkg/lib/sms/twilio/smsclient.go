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
	Client *http.Client

	AccountSID string

	AuthToken    string
	APIKey       string
	APIKeySecret string

	From                string
	MessagingServiceSID string

	Logger *slog.Logger
}

func (t *TwilioClient) send(ctx context.Context, options *smsclient.SendOptions) ([]byte, *SendResponse, error) {
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

	// https://www.twilio.com/docs/usage/api#authenticate-with-http
	if t.AuthToken != "" {
		// When Auth Token is used, username is Account SID, and password is Auth Token.
		req.SetBasicAuth(t.AccountSID, t.AuthToken)
	} else if t.APIKey != "" {
		// When API Key is used, username is API key, and password is API key secret.
		req.SetBasicAuth(t.APIKey, t.APIKeySecret)
	} else { //nolint: staticcheck
		// Normally we should not reach here.
		// But in case we do, we do not provide the auth header.
		// And Twilio should returns an error response to us in this case.
	}

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
			&smsclient.SendResultError{
				DumpedResponse: dumpedResponse,
			},
		)
	}

	sendResponse, err := ParseSendResponse(respData)
	if err != nil {
		return nil, nil, errors.Join(
			err,
			&smsclient.SendResultError{
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

	t.Logger.LogAttrs(ctx, slog.LevelInfo, "twilio response", attrs...)

	return dumpedResponse, sendResponse, nil
}

func (t *TwilioClient) Send(ctx context.Context, options *smsclient.SendOptions) (*smsclient.SendResultSuccess, error) {
	ctx = smsclient.WithSendContext(ctx, func(sendCtx *smsclient.SendContext) {
		sendCtx.Twilio = &smsclient.SendContextTwilio{
			BodyLength: len(options.Body),
		}
	})

	dumpedResponse, sendSMSResponse, err := t.send(ctx, options)
	if err != nil {
		return nil, err
	}

	var segmentCount *int
	if sendSMSResponse.NumSegments != nil {
		if parsed, err := strconv.Atoi(*sendSMSResponse.NumSegments); err == nil {
			segmentCount = &parsed
		}
	}

	_ = smsclient.WithSendContext(ctx, func(sendCtx *smsclient.SendContext) {
		sendCtx.Twilio.SegmentCount = segmentCount
	})

	// Success case.
	if sendSMSResponse.ErrorCode == nil {
		return &smsclient.SendResultSuccess{
			DumpedResponse: dumpedResponse,
		}, nil
	}

	// Failed case.
	return nil, &smsclient.SendResultError{
		DumpedResponse: dumpedResponse,
	}
}

var _ smsclient.RawClient = &TwilioClient{}
