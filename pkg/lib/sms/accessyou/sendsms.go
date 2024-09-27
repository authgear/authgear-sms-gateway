package accessyou

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
)

//go:generate mockgen -source=sendsms.go -destination=sendsms_mock_test.go -package accessyou

var leadingBOMRegexp = regexp.MustCompile(`^[\x{feff}]+`)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func fixRespData(respData []byte) []byte {
	// The response data is in format
	// "\ufeff{\"msg_status\":\"100\",\"msg_status_desc\":\"Successfully submitted message. \\u6267\\u884c\\u6210\\u529f\",\"phoneno\":\"852********\",\"msg_id\":852309279}"
	// Remove BOM token from resp json
	return leadingBOMRegexp.ReplaceAll(respData, []byte(""))
}

func SendSMS(
	client HTTPClient,
	baseUrl string,
	accountNo string,
	user string,
	pwd string,
	sender string,
	to string,
	body string,
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
