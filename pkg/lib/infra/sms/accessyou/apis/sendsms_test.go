package apis

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms/accessyou/models"
)

var successResponse = []byte(
	"\ufeff{\"msg_status\":\"100\",\"msg_status_desc\":\"Successfully submitted message. \\u6267\\u884c\\u6210\\u529f\",\"phoneno\":\"85264975244\",\"msg_id\":854998103}",
)

func TestSendSMS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := NewMockHTTPClient(ctrl)

	test := func(
		baseUrl string,
		accountNo string,
		user string,
		pwd string,
		sender string,
		to string,
		body string,
		expect func(clientDo *gomock.Call),
		callback func([]byte, *models.SendSMSResponse, error),
	) {
		u, err := url.Parse(baseUrl)
		if err != nil {
			panic(err)
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

		expect(client.EXPECT().Do(req))

		respData, sendSMSResponse, err := SendSMS(
			client,
			baseUrl,
			accountNo,
			user,
			pwd,
			sender,
			to,
			body,
		)

		callback(respData, sendSMSResponse, err)
	}

	Convey("SendSMS success", t, func() {
		var baseUrl = "https://www.example.com"
		var accountNo = "accountno"
		var pwd = "pwd"
		var to = "to"
		var body = "This is your OTP 123456"
		var user = "user"
		var sender = "sender"
		test(
			baseUrl, accountNo, user, pwd, sender, to, body,
			func(clientDo *gomock.Call) {
				clientDo.Return(
					&http.Response{
						Body: io.NopCloser(bytes.NewReader(successResponse)),
					},
					nil)
			},
			func(respData []byte, sendSMSResponse *models.SendSMSResponse, err error) {
				So(err, ShouldBeNil)
				So(respData, ShouldEqual, []byte(
					"{\"msg_status\":\"100\",\"msg_status_desc\":\"Successfully submitted message. \\u6267\\u884c\\u6210\\u529f\",\"phoneno\":\"85264975244\",\"msg_id\":854998103}",
				))
				So(sendSMSResponse.Status, ShouldEqual, "100")
			},
		)
	})
}
