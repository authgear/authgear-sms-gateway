package sms

import (
	"fmt"
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
) (*smsclient.SendResult, error) {
	clientName := GetClientNameByMatch(s.RootConfig, &MatchContext{AppID: appID, PhoneNumber: string(sendOptions.To)})
	client := s.SMSClientMap.GetClientByName(clientName)
	s.Logger.Info(fmt.Sprintf("Client %v is selected for %v", clientName, sendOptions.To))
	return client.Send(sendOptions)
}
