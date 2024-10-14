package sms

import (
	"context"
	"errors"
	"log/slog"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/logger"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type SMSService struct {
	Logger       *slog.Logger
	RootConfig   *config.RootConfig
	SMSClientMap SMSClientMap
}

func (s *SMSService) Send(
	ctx context.Context,
	appID string,
	sendOptions *smsclient.SendOptions,
) (*smsclient.SendResult, error) {
	clientName := GetClientNameByMatch(s.RootConfig, &MatchContext{AppID: appID, PhoneNumber: string(sendOptions.To)})
	client := s.SMSClientMap.GetClientByName(clientName)

	ctx = logger.ContextWithAttrs(ctx, slog.String("client_name", clientName))

	s.Logger.InfoContext(ctx, "selected client")

	result, err := client.Send(ctx, sendOptions)
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
