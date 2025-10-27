package sendcloud

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
	"strings"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/api"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sensitive"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

func Send(ctx context.Context, client *http.Client, baseUrl string, sendRequest *SendRequest, smsKey string, logger *slog.Logger) ([]byte, *SendResponse, error) {
	values := sendRequest.ToValues()
	values.Set("signature", sendRequest.Sign(smsKey))

	data := values.Encode()

	req, _ := http.NewRequest("POST", fmt.Sprintf("%v/smsapi/send", baseUrl), strings.NewReader(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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
			&smsclient.SendResultError{
				DumpedResponse: dumpedResponse,
			},
		)
	}

	sendResponse, err := ParseSendResponse(respData)
	if err != nil {
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

	logger.InfoContext(ctx, "sendcloud response",
		"result", sendResponse.Result,
		"statusCode", sendResponse.StatusCode,
		"message", sendResponse.Message,
	)

	return dumpedResponse, sendResponse, nil
}
