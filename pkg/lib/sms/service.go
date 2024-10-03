package sms

import (
	"errors"
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
	s.Logger.Info("selected client",
		"to", sendOptions.To,
		"client_name", clientName,
	)

	result, err := client.Send(sendOptions)
	var errSendResult *smsclient.SendResult
	if errors.As(err, &errSendResult) {
		if errSendResult.Info == nil {
			errSendResult.Info = &smsclient.SendResultInfo{}
		}
		if errSendResult.Info.SendResultInfoRoot == nil {
			errSendResult.Info.SendResultInfoRoot = &smsclient.SendResultInfoRoot{}
		}

		errSendResult.Info.SendResultInfoRoot.ProviderName = clientName
	}
	if err != nil {
		return nil, err
	}

	if result.Info == nil {
		result.Info = &smsclient.SendResultInfo{}
	}
	if result.Info.SendResultInfoRoot == nil {
		result.Info.SendResultInfoRoot = &smsclient.SendResultInfoRoot{}
	}
	result.Info.SendResultInfoRoot = &smsclient.SendResultInfoRoot{
		ProviderName: clientName,
	}

	return result, nil
}
