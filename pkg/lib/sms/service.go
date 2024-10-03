package sms

import (
	"log/slog"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type SMSService struct {
	Logger       *slog.Logger
	RootConfig   *config.RootConfig
	SMSClientMap SMSClientMap
}

func (s *SMSService) Send(
	appID string,
	sendOptions *smsclient.SendOptions,
) (*smsclient.SendResult, *smsclient.SendResultInfo, error) {
	clientName := GetClientNameByMatch(s.RootConfig, &MatchContext{AppID: appID, PhoneNumber: string(sendOptions.To)})
	client := s.SMSClientMap.GetClientByName(clientName)
	s.Logger.Info("selected client",
		"to", sendOptions.To,
		"client_name", clientName,
	)

	result, info, err := client.Send(sendOptions)

	info.SendResultInfoRoot = &smsclient.SendResultInfoRoot{
		ProviderName: clientName,
	}

	return result, info, err
}
