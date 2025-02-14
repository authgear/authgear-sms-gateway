package accessyou

import (
	"context"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/api"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type AccessYouClient struct {
	BaseUrl   string
	Client    *http.Client
	AccountNo string
	User      string
	Pwd       string
	From      string
	Logger    *slog.Logger
}

func NewAccessYouClient(
	httpClient *http.Client,
	baseUrl string,
	accountNo string,
	user string,
	pwd string,
	from string,
	logger *slog.Logger,
) *AccessYouClient {
	if baseUrl == "" {
		baseUrl = "http://sms.accessyou-anyip.com"
	}
	return &AccessYouClient{
		Client:    httpClient,
		BaseUrl:   baseUrl,
		AccountNo: accountNo,
		User:      user,
		Pwd:       pwd,
		From:      from,
		Logger:    logger,
	}
}

var plusHyphensRegexp = regexp.MustCompile(`[\+\-]+`)

func fixPhoneNumber(phoneNumber string) string {
	// Access you phone number should have no + and -
	return plusHyphensRegexp.ReplaceAllString(phoneNumber, "")
}

func (n *AccessYouClient) Send(ctx context.Context, options *smsclient.SendOptions) (*smsclient.SendResultSuccess, error) {
	to := fixPhoneNumber(string(options.To))

	dumpedResponse, sendSMSResponse, err := SendSMS(
		ctx,
		n.Client,
		n.BaseUrl,
		n.AccountNo,
		n.User,
		n.Pwd,
		n.From,
		to,
		options.Body,
		n.Logger,
	)
	if err != nil {
		return nil, err
	}

	// Success case.
	if sendSMSResponse.Status == "100" {
		return &smsclient.SendResultSuccess{
			DumpedResponse: dumpedResponse,
		}, nil
	}

	// Failed case.
	return nil, n.makeError(sendSMSResponse.Status, dumpedResponse)
}

func (n *AccessYouClient) makeError(
	msgStatus string,
	dumpedResponse []byte,
) *smsclient.SendResultError {
	err := &smsclient.SendResultError{
		DumpedResponse:    dumpedResponse,
		ProviderName:      "accessyou",
		ProviderErrorCode: msgStatus,
	}

	// See https://www.accessyou.com/smsapi.pdf
	switch msgStatus {
	case "108":
		fallthrough
	case "110":
		err.Code = api.CodeInvalidPhoneNumber
	case "105":
		err.Code = api.CodeAuthenticationFailed
	case "106":
		fallthrough
	case "107":
		fallthrough
	case "too_many_login_failure":
		err.Code = api.CodeDeliveryRejected
	}

	return err
}

var _ smsclient.RawClient = &AccessYouClient{}
