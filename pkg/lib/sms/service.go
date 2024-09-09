package sms

import (
	"fmt"
	"log/slog"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/type_util"
)

type SMSService struct {
	Logger              *slog.Logger
	SMSProviderSelector *SMSProviderSelector
}

func NewSMSService(logger *slog.Logger, smsProviderConfig *config.SMSProviderConfig) (*SMSService, error) {
	smsProviders, err := NewSMSProviders(smsProviderConfig, logger)
	if err != nil {
		return nil, err
	}
	smsProviderSelector, err := NewSMSProviderSelector(smsProviderConfig, smsProviders)
	if err != nil {
		return nil, err
	}
	return &SMSService{
		Logger:              logger,
		SMSProviderSelector: smsProviderSelector,
	}, nil
}

func (s *SMSService) Send(to type_util.SensitivePhoneNumber, body string) error {
	client, err := s.SMSProviderSelector.GetClientByMatch(&MatchContext{PhoneNumber: string(to)})
	if err != nil {
		return err
	}
	s.Logger.Info(fmt.Sprintf("Client %v is selected for %v", client.GetName(), to))
	err = client.Send(string(to), body)
	return nil
}
