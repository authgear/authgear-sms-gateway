package sms

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms/sendcloud"
)

type SendOptions struct {
	To                string
	Body              string
	TemplateName      string
	LanguageTag       string
	TemplateVariables *TemplateVariables
}

type SendResult struct {
	ClientResponse []byte
	Success        bool
}

type RawClient interface {
	Send(options *SendOptions) (*SendResult, error)
}

func NewClientFromConfigProvider(p *config.Provider, logger *slog.Logger) RawClient {
	switch p.Type {
	case config.ProviderTypeTwilio:
		return NewTwilioClient(
			p.Twilio.AccountSID,
			p.Twilio.AuthToken,
			p.Twilio.Sender,
			p.Twilio.MessagingServiceSID,
		)
	case config.ProviderTypeNexmo:
		return NewNexmoClient(
			p.Nexmo.APIKey,
			p.Nexmo.APISecret,
			p.Nexmo.Sender,
		)
	case config.ProviderTypeAccessYou:
		return NewAccessYouClient(
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
		return NewSendCloudClient(
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
