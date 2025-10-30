package accessyouotp

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/api"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/accessyou"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type AccessYouOTPClient struct {
	BaseUrl   string
	Client    *http.Client
	AccountNo string
	User      string
	Pwd       string
	TID       string
	Logger    *slog.Logger
}

func NewAccessYouOTPClient(
	httpClient *http.Client,
	baseUrl string,
	accountNo string,
	user string,
	pwd string,
	tid string,
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
		TID:       tid,
		Logger:    logger,
	}
}

func (n *AccessYouOTPClient) Send(ctx context.Context, options *smsclient.SendOptions) (*smsclient.SendResultSuccess, error) {
	to := accessyou.FixPhoneNumber(string(options.To))

	code := options.TemplateVariables.Code
	if code == "" {
		return nil, &smsclient.SendResultError{
			Code:              api.CodeUnsupportedRequest,
			Description:       "Only 6-digit OTP message is supported.",
			DumpedResponse:    nil,
			ProviderErrorCode: "",
		}
	}

	dumpedResponse, sendSMSResponse, err := SendOTPSMS(
		ctx,
		n.Client,
		n.BaseUrl,
		n.Logger,
		&SendOTPSMSOptions{
			AccountNo: n.AccountNo,
			User:      n.User,
			Pwd:       n.Pwd,
			TID:       n.TID,
			To:        to,
			AppName:   options.TemplateVariables.AppName,
			Code:      code,
		},
	)
	if err != nil {
		return nil, err
	}

	// Success case.
	if sendSMSResponse.Status == accessyou.STATUS_SUCCESS {
		return &smsclient.SendResultSuccess{
			DumpedResponse: dumpedResponse,
		}, nil
	}

	// Failed case.
	return nil, accessyou.MakeError(sendSMSResponse.Status, dumpedResponse)
}

func (n *AccessYouOTPClient) ProviderType() string {
	return "accessyou_otp"
}

var _ smsclient.RawClient = &AccessYouOTPClient{}
