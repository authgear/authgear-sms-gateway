package accessyou

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
)

var leadingBOMRegexp = regexp.MustCompile(`^[\x{feff}]+`)

func fixRespData(respData []byte) []byte {
	// Remove BOM token from resp json,
	// See _test.go for details.
	return leadingBOMRegexp.ReplaceAll(respData, []byte(""))
}

func SendSMS(
	client *http.Client,
	baseUrl string,
	accountNo string,
	user string,
	pwd string,
	sender string,
	to string,
	body string,
) ([]byte, *SendSMSResponse, error) {
	// TODO: Add logs
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
		"POST",
		u.String(),
		nil)
	req.Header.Set("Cookie", "dynamic=sms")

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	respData = fixRespData(respData)

	sendSMSResponse, err := ParseSendSMSResponse(respData)
	if err != nil {
		return respData, nil, err
	}

	return respData, sendSMSResponse, err
}
