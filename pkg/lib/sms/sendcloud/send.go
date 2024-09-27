package sendcloud

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Send(client *http.Client, baseUrl string, sendRequest *SendRequest, smsKey string) ([]byte, *SendResponse, error) {
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

	respData, err := io.ReadAll(resp.Body)

	sendResponse, err := ParseSendResponse(respData)
	if err != nil {
		return respData, nil, err
	}

	return respData, sendResponse, err
}
