package accessyouotp

import (
	"context"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/accessyou"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type AccessYouOTPClient struct {
	BaseUrl   string
	Client    *http.Client
	AccountNo string
	User      string
	Pwd       string
	A         string
	Logger    *slog.Logger
}

func NewAccessYouOTPClient(
	httpClient *http.Client,
	baseUrl string,
	accountNo string,
	user string,
	pwd string,
	a string,
	logger *slog.Logger,
) *AccessYouOTPClient {
	if baseUrl == "" {
		baseUrl = "https://otp.accessyou-anyip.com"
	}
	return &AccessYouOTPClient{
		Client:    httpClient,
		BaseUrl:   baseUrl,
		AccountNo: accountNo,
		User:      user,
		Pwd:       pwd,
		A:         a,
		Logger:    logger,
	}
}

var plusHyphensRegexp = regexp.MustCompile(`[\+\-]+`)

func fixPhoneNumber(phoneNumber string) string {
	// Access you phone number should have no + and -
	return plusHyphensRegexp.ReplaceAllString(phoneNumber, "")
}

func (n *AccessYouOTPClient) Send(ctx context.Context, options *smsclient.SendOptions) (*smsclient.SendResultSuccess, error) {
	to := fixPhoneNumber(string(options.To))

	dumpedResponse, sendSMSResponse, err := SendOTPSMS(
		ctx,
		n.Client,
		n.BaseUrl,
		n.Logger,
		&SendOTPSMSOptions{
			AccountNo: n.AccountNo,
			User:      n.User,
			Pwd:       n.Pwd,
			A:         n.A,
			To:        to,
			Code:      options.Body,
		},
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
	return nil, accessyou.MakeError(sendSMSResponse.Status, dumpedResponse, "accessyouotp")
}

var _ smsclient.RawClient = &AccessYouOTPClient{}
