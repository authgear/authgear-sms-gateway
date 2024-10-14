package accessyou

import (
	"context"
	"log/slog"
	"net/http"
	"regexp"

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
	return nil, &smsclient.SendResultError{
		DumpedResponse: dumpedResponse,
	}
}

var _ smsclient.RawClient = &AccessYouClient{}
