package accessyou

import (
	"bytes"
	"net/http"
	"testing"

	"gopkg.in/h2non/gock.v1"

	. "github.com/smartystreets/goconvey/convey"
)

var successResponseWithoutBOM = `{"msg_status":"100","msg_status_desc":"Successfully submitted message. 执行成功","phoneno":"85264975244","msg_id":854998103}`

var successResponseWithBOM = "\ufeff" + successResponseWithoutBOM

func TestSendSMS(t *testing.T) {
	Convey("SendSMS success", t, func() {
		var baseUrl = "https://www.example.com"
		var accountNo = "accountno"
		var pwd = "pwd"
		var to = "to"
		var body = "This is your OTP 123456"
		var user = "user"
		var sender = "sender"

		httpClient := &http.Client{}
		gock.InterceptClient(httpClient)
		defer gock.Off()

		gock.New("https://www.example.com").
			Post("/sendsms.php").
			Reply(200).
			Body(bytes.NewReader([]byte(successResponseWithBOM)))

		rawBody, parsedResponse, err := SendSMS(httpClient, baseUrl, accountNo, user, pwd, sender, to, body)

		So(err, ShouldBeNil)
		So(rawBody, ShouldResemble, []byte(successResponseWithoutBOM))
		So(parsedResponse, ShouldResemble, &SendSMSResponse{
			MessageID:   854998103,
			Status:      "100",
			Description: "Successfully submitted message. 执行成功",
			PhoneNo:     "85264975244",
		})
	})
}
