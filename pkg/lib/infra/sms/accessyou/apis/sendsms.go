package apis

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms/accessyou/models"
)

var leadingBOMRegexp = regexp.MustCompile(`^[\x{feff}]+`)

func fixRespData(respData []byte) []byte {
	// The response data is in format
	// "\ufeff{\"msg_status\":\"100\",\"msg_status_desc\":\"Successfully submitted message. \\u6267\\u884c\\u6210\\u529f\",\"phoneno\":\"852********\",\"msg_id\":852309279}"
	// Remove BOM token from resp json
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
) ([]byte, *models.SendSMSResponse, error) {
	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"%v/sendsms.php?accountno=%v&pwd=%v&tid=1&phone=%v&a=%v&user=%v&from=%v",
			baseUrl,
			accountNo,
			pwd,
			to,
			body,
			user,
			sender),
		nil)
	req.Header.Set("Cookie", "dynamic=sms")

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)

	respData = fixRespData(respData)

	sendSMSResponse, err := models.ParseSendSMSResponse(respData)
	if err != nil {
		return respData, nil, err
	}
	return respData, sendSMSResponse, err
}
