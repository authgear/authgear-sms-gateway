package accessyou

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sensitive"
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
		return nil, nil, sensitive.RedactHTTPClientError(err)
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
		return nil, nil, errors.Join(
			err,
			MakeError("", dumpedResponse),
		)
	}

	logger.InfoContext(ctx, "accessyou response",
		"msg_id", sendSMSResponse.MessageID,
		"msg_status", sendSMSResponse.Status,
		"msg_status_desc", sendSMSResponse.Description,
	)

	return dumpedResponse, sendSMSResponse, nil
}
