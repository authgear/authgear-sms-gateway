package sendcloud

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

func Send(client *http.Client, baseUrl string, sendRequest *SendRequest, smsKey string, logger *slog.Logger) ([]byte, *SendResponse, error) {
	values := sendRequest.ToValues()
	values.Set("signature", sendRequest.Sign(smsKey))

	data := values.Encode()

	req, _ := http.NewRequest("POST", fmt.Sprintf("%v/smsapi/send", baseUrl), strings.NewReader(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

	logger.Info("sendcloud response",
		"result", sendResponse.Result,
		"statusCode", sendResponse.StatusCode,
		"message", sendResponse.Message,
	)

	return dumpedResponse, sendResponse, nil
}
