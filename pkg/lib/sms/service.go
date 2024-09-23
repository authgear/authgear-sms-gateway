package sms

import (
	"fmt"
	"log/slog"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	. "github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/type_util"
)

type SMSService struct {
	Logger            *slog.Logger
	SMSProviderConfig *config.SMSProviderConfig
	SMSClientMap      SMSClientMap
}

func NewSMSService(
	logger *slog.Logger,
	smsProviderConfig *config.SMSProviderConfig,
	smsClientMap SMSClientMap,
) *SMSService {
	return &SMSService{
		Logger:            logger,
		SMSProviderConfig: smsProviderConfig,
		SMSClientMap:      smsClientMap,
	}
}

func (s *SMSService) Send(
	appID string,
	to type_util.SensitivePhoneNumber,
	body string,
	templateName string,
	languageTag string,
	templateVariables *TemplateVariables,
) (*SendResult, error) {
	clientName := GetClientNameByMatch(s.SMSProviderConfig, &MatchContext{AppID: appID, PhoneNumber: string(to)})
	client := s.SMSClientMap.GetClientByName(clientName)
	s.Logger.Info(fmt.Sprintf("Client %v is selected for %v", clientName, to))
	return client.Send(&SendOptions{
		To:                string(to),
		Body:              body,
		TemplateName:      templateName,
		LanguageTag:       languageTag,
		TemplateVariables: templateVariables,
	})
}
