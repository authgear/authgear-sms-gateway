package accessyouotp

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/accessyou"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type SendOTPSMSOptions struct {
	AccountNo string
	User      string
	Pwd       string
	A         string
	To        string
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
		"tid":       {"1"},
		"phone":     {opts.To},
		"a":         {opts.A},
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

	respData = accessyou.FixRespData(respData)
	sendSMSResponse, err := accessyou.ParseSendSMSResponse(respData)
	if err != nil {
		return nil, nil, errors.Join(
			err,
			&smsclient.SendResultError{
				DumpedResponse: dumpedResponse,
			},
		)
	}

	logger.InfoContext(ctx, "accessyouotp response",
		"msg_id", sendSMSResponse.MessageID,
		"msg_status", sendSMSResponse.Status,
		"msg_status_desc", sendSMSResponse.Description,
	)

	return dumpedResponse, sendSMSResponse, nil
}
