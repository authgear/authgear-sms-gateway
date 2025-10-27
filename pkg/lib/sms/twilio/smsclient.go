package twilio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/api"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sensitive"
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
		err = sensitive.RedactHTTPClientError(err)
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			err = errors.Join(err, &smsclient.SendResultError{
				DumpedResponse: nil,
				Code:           api.CodeProviderTimeout,
			})
		}
		return nil, nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

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
		var jsonUnmarshalErr *json.UnmarshalTypeError
		if errors.As(err, &jsonUnmarshalErr) {
			return nil, nil, t.parseAndHandleErrorResponse(respData, dumpedResponse)
		}
		sendErr := &smsclient.SendResultError{
			DumpedResponse: dumpedResponse,
		}
		var jsonSyntaxErr *json.SyntaxError
		if errors.As(err, &jsonSyntaxErr) {
			sendErr.Code = api.CodeUnknownResponseFormat
		}
		return nil, nil, errors.Join(
			err,
			sendErr,
		)
	}

	if sendResponse.ErrorCode != nil {
		return nil, nil, t.makeError(*sendResponse.ErrorCode, dumpedResponse)
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

	return &smsclient.SendResultSuccess{
		DumpedResponse: dumpedResponse,
	}, nil

}

func (t *TwilioClient) parseAndHandleErrorResponse(
	responseBody []byte,
	dumpedResponse []byte,
) error {
	errResponse, err := ParseErrorResponse(responseBody)

	if err != nil {
		var jsonUnmarshalErr *json.UnmarshalTypeError
		if errors.As(err, &jsonUnmarshalErr) {
			// Not something we can understand, return an error with the dumped response
			return &smsclient.SendResultError{
				DumpedResponse: dumpedResponse,
			}
		} else {
			return errors.Join(err, &smsclient.SendResultError{
				DumpedResponse: dumpedResponse,
			})
		}
	}

	return t.makeError(errResponse.Code, dumpedResponse)
}

func (t *TwilioClient) makeError(
	errorCode int,
	dumpedResponse []byte,
) *smsclient.SendResultError {
	err := &smsclient.SendResultError{
		DumpedResponse:    dumpedResponse,
		ProviderErrorCode: fmt.Sprintf("%d", errorCode),
	}

	// See https://www.twilio.com/docs/api/errors
	switch errorCode {
	case 21211: // Invalid 'To' Phone Number
		fallthrough
	case 21265: // 'To' number cannot be a Short Code
		err.Code = api.CodeInvalidPhoneNumber
	case 30022:
		fallthrough
	case 14107:
		fallthrough
	case 51002:
		fallthrough
	case 63017:
		fallthrough
	case 63018:
		err.Code = api.CodeRateLimited
	case 20003:
		err.Code = api.CodeAuthenticationFailed
	case 30002: // Account suspended
		fallthrough
	case 21264: // ‘From’ phone number not verified
		fallthrough
	case 21266: // 'To' and 'From' numbers cannot be the same
		fallthrough
	case 21267: // Alphanumeric Sender ID cannot be used as the 'From' number on trial accounts
		fallthrough
	case 21606: // The 'From' phone number provided is not a valid message-capable Twilio phone number for this destination/account
		fallthrough
	case 21607: // The 'from' phone number must be the sandbox phone number for trial accounts.
		fallthrough
	case 21659: // 'From' is not a Twilio phone number or Short Code country mismatch
		fallthrough
	case 21660: // Mismatch between the 'From' number and the account
		fallthrough
	case 21661: // 'From' number is not SMS-capable
		fallthrough
	case 21910: // Invalid 'From' and 'To' pair. 'From' and 'To' should be of the same channel
		fallthrough
	case 63007: // Twilio could not find a Channel with the specified 'From' address
		err.Code = api.CodeDeliveryRejected
	}

	return err
}

func (t *TwilioClient) ProviderType() string {
	return "twilio"
}

var _ smsclient.RawClient = &TwilioClient{}
