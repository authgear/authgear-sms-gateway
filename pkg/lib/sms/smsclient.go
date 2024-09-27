package sms

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/accessyou"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/sendcloud"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/twilio"
)

func NewClientFromConfigProvider(p *config.Provider, httpClient *http.Client, logger *slog.Logger) smsclient.RawClient {
	switch p.Type {
	case config.ProviderTypeTwilio:
		return twilio.NewTwilioClient(
			p.Twilio.AccountSID,
			p.Twilio.AuthToken,
			p.Twilio.Sender,
			p.Twilio.MessagingServiceSID,
		)
	case config.ProviderTypeAccessYou:
		return accessyou.NewAccessYouClient(
			httpClient,
			p.AccessYou.BaseUrl,
			p.AccessYou.AccountNo,
			p.AccessYou.User,
			p.AccessYou.Pwd,
			p.AccessYou.Sender,
			logger,
		)
	case config.ProviderTypeSendCloud:
		templateResolver := sendcloud.NewSendCloudTemplateResolver(
			p.SendCloud.Templates,
			p.SendCloud.TemplateAssignments,
		)
		return sendcloud.NewSendCloudClient(
			httpClient,
			p.SendCloud.BaseUrl,
			p.SendCloud.SMSUser,
			p.SendCloud.SMSKey,
			templateResolver,
			logger,
		)
	default:
		panic(errors.New(fmt.Sprintf("Unknown type %s", p.Type)))
	}
}
