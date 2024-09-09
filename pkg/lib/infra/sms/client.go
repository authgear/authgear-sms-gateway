package sms

import (
	"errors"
	"fmt"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
)

type RawClient interface {
	Send(to string, body string) error
	GetName() string
}

func NewClientFromConfigProvider(p *config.Provider) (RawClient, error) {
	switch p.Type {
	case config.ProviderTypeTwilio:
		return NewTwilioClient(
			p.Name,
			p.Twilio.AccountSID,
			p.Twilio.AuthToken,
			p.Twilio.Sender,
			p.Twilio.MessagingServiceSID,
		), nil
	case config.ProviderTypeNexmo:
		return NewNexmoClient(
			p.Name,
			p.Nexmo.APIKey,
			p.Nexmo.APISecret,
			p.Nexmo.Sender,
		), nil
	case config.ProviderTypeAccessYou:
		fallthrough
	case config.ProviderTypeSendCloud:
		fallthrough
	case config.ProviderTypeInfobip:
		fallthrough
	default:
		return nil, errors.New(fmt.Sprintf("Unknown type %s", p.Type))
	}
}
