package sms

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms/sendcloud"
)

type ClientResponse []byte

type RawClient interface {
	Send(to string, body string, templateName string, languageTag string, templateVariables *TemplateVariables) (ClientResponse, error)
	GetName() string
}

func NewClientFromConfigProvider(p *config.Provider, logger *slog.Logger) RawClient {
	switch p.Type {
	case config.ProviderTypeTwilio:
		return NewTwilioClient(
			p.Name,
			p.Twilio.AccountSID,
			p.Twilio.AuthToken,
			p.Twilio.Sender,
			p.Twilio.MessagingServiceSID,
		)
	case config.ProviderTypeNexmo:
		return NewNexmoClient(
			p.Name,
			p.Nexmo.APIKey,
			p.Nexmo.APISecret,
			p.Nexmo.Sender,
		)
	case config.ProviderTypeAccessYou:
		return NewAccessYouClient(
			p.Name,
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
			p.Name,
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
