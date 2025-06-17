package accessyouotp

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"gopkg.in/h2non/gock.v1"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/accessyou"
	. "github.com/smartystreets/goconvey/convey"
)

var successResponseWithoutBOM = `{"msg_status":"100","msg_status_desc":"Successfully submitted message. 执行成功","phoneno":"85264975244","msg_id":854998103}`

var successResponseWithBOM = "\ufeff" + successResponseWithoutBOM

func TestSendOTPSMS(t *testing.T) {
	Convey("SendOTPSMS success", t, func() {
		var baseUrl = "https://www.example.com"
		var accountNo = "accountno"
		var pwd = "pwd"
		var to = "to"
		var tid = "1"
		var appName = "appName"
		var code = "123456"
		var user = "user"

		httpClient := &http.Client{}
		gock.InterceptClient(httpClient)
		defer gock.Off()

		gock.New("https://www.example.com").
			Get("/sendsms-otp.php").
			MatchParam("accountno", accountNo).
			MatchParam("pwd", pwd).
			MatchParam("tid", tid).
			MatchParam("phone", to).
			MatchParam("a", appName).
			MatchParam("b", code).
			MatchParam("user", user).
			Reply(200).
			BodyString(successResponseWithBOM)

		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		ctx := context.Background()
		dumpedResponse, parsedResponse, err := SendOTPSMS(
			ctx,
			httpClient,
			baseUrl,
			logger,
			&SendOTPSMSOptions{
				AccountNo: accountNo,
				User:      user,
				Pwd:       pwd,
				TID:       tid,
				To:        to,
				AppName:   appName,
				Code:      code,
			},
		)

		So(err, ShouldBeNil)
		So(string(dumpedResponse), ShouldEqual, "HTTP/1.1 200 OK\r\nContent-Length: 131\r\n\r\n"+successResponseWithBOM)
		So(parsedResponse, ShouldResemble, &accessyou.SendSMSResponse{
			MessageID:   854998103,
			Status:      "100",
			Description: "Successfully submitted message. 执行成功",
			PhoneNo:     "85264975244",
		})
	})
}
