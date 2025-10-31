package accessyouotp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/api"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sensitive"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/accessyou"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type SendOTPSMSOptions struct {
	AccountNo string
	User      string
	Pwd       string
	TID       string
	To        string
	AppName   string
	Code      string
}

func SendOTPSMS(
	ctx context.Context,
	client *http.Client,
	baseUrl string,
	logger *slog.Logger,
	opts *SendOTPSMSOptions,
) ([]byte, *accessyou.SendSMSResponse, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, nil, err
	}
	u.Path = "/sendsms-otp.php"

	queryParams := url.Values{
		"accountno": {opts.AccountNo},
		"pwd":       {opts.Pwd},
		"tid":       {opts.TID},
		"phone":     {opts.To},
		"a":         {opts.AppName},
		"b":         {opts.Code},
		"user":      {opts.User},
	}
	u.RawQuery = queryParams.Encode()

	req, _ := http.NewRequest(
		"GET",
		u.String(),
		nil)
	req.Header.Set("Cookie", "dynamic=otp")

	resp, err := client.Do(req)
	if err != nil {
		err = sensitive.RedactHTTPClientError(err)
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			err = errors.Join(err, &smsclient.SendResultError{
				DumpedResponse: nil,
				Code:           api.CodeTimeout,
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
			accessyou.MakeError("", dumpedResponse),
		)
	}

	respData = accessyou.FixRespData(respData)
	sendSMSResponse, err := accessyou.ParseSendSMSResponse(respData)
	if err != nil {
		sendErr := &smsclient.SendResultError{
			DumpedResponse: dumpedResponse,
		}
		var jsonSyntaxErr *json.SyntaxError
		if errors.As(err, &jsonSyntaxErr) {
			sendErr.Code = api.CodeUnknownError
		}
		return nil, nil, errors.Join(
			err,
			sendErr,
		)
	}

	logger.InfoContext(ctx, "accessyou_otp response",
		"msg_id", sendSMSResponse.MessageID,
		"msg_status", sendSMSResponse.Status,
		"msg_status_desc", sendSMSResponse.Description,
	)

	return dumpedResponse, sendSMSResponse, nil
}
