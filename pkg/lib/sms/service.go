package sms

import (
	"context"
	"log/slog"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/logger"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type SMSService struct {
	Logger         *slog.Logger
	RootConfig     *config.RootConfig
	SMSProviderMap SMSProviderMap
}

func (s *SMSService) Send(
	ctx context.Context,
	appID string,
	sendOptions *smsclient.SendOptions,
) (*smsclient.SendResultSuccess, error) {
	clientName := GetProviderNameByMatch(s.RootConfig, &MatchContext{AppID: appID, PhoneNumber: string(sendOptions.To)})
	client := s.SMSProviderMap.GetProviderByName(clientName)

	ctx = smsclient.WithSendContext(ctx, func(sendCtx *smsclient.SendContext) {
		if sendCtx.Root == nil {
			sendCtx.Root = &smsclient.SendContextRoot{}
		}
		sendCtx.Root.ProviderName = clientName
	})
	ctx = logger.ContextWithAttrs(ctx, slog.String("provider_name", clientName))

	s.Logger.InfoContext(ctx, "selected provider")

	result, err := client.Send(ctx, sendOptions)
	if err != nil {
		return nil, err
	}

	return result, nil
}
