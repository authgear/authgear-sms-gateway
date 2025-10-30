package accessyou

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
	"regexp"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/api"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sensitive"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

var leadingBOMRegexp = regexp.MustCompile(`^[\x{feff}]+`)

func FixRespData(respData []byte) []byte {
	// Remove BOM token from resp json,
	// See _test.go for details.
	return leadingBOMRegexp.ReplaceAll(respData, []byte(""))
}

func SendSMS(
	ctx context.Context,
	client *http.Client,
	baseUrl string,
	accountNo string,
	user string,
	pwd string,
	sender string,
	to string,
	body string,
	logger *slog.Logger,
) ([]byte, *SendSMSResponse, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, nil, err
	}
	u.Path = "/sendsms.php"

	queryParams := url.Values{
		"accountno": {accountNo},
		"pwd":       {pwd},
		"tid":       {"1"},
		"phone":     {to},
		"a":         {body},
		"user":      {user},
		"from":      {sender},
	}
	u.RawQuery = queryParams.Encode()

	req, _ := http.NewRequest(
		"GET",
		u.String(),
		nil)
	req.Header.Set("Cookie", "dynamic=sms")

	resp, err := client.Do(req)
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
			MakeError("", dumpedResponse),
		)
	}

	respData = FixRespData(respData)
	sendSMSResponse, err := ParseSendSMSResponse(respData)
	if err != nil {
		var jsonSyntaxErr *json.SyntaxError
		sendErr := &smsclient.SendResultError{
			DumpedResponse: dumpedResponse,
		}
		if errors.As(err, &jsonSyntaxErr) {
			sendErr.Code = api.CodeUnknownResponseFormat
		}
		return nil, nil, errors.Join(
			err,
			sendErr,
		)
	}

	logger.InfoContext(ctx, "accessyou response",
		"msg_id", sendSMSResponse.MessageID,
		"msg_status", sendSMSResponse.Status,
		"msg_status_desc", sendSMSResponse.Description,
	)

	return dumpedResponse, sendSMSResponse, nil
}
