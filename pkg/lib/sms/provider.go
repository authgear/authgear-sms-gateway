package sms

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms"
)

type SMSClientMap map[string]sms.RawClient

func NewSMSClientMap(c *config.SMSProviderConfig, logger *slog.Logger) SMSClientMap {
	var clientMap = make(map[string]sms.RawClient)

	for _, provider := range c.Providers {
		client := sms.NewClientFromConfigProvider(provider, logger)
		clientMap[provider.Name] = client
	}

	return SMSClientMap(clientMap)
}

func (s SMSClientMap) GetClientByName(name string) sms.RawClient {
	client := s[name]
	if client == nil {
		panic(errors.New(fmt.Sprintf("Unknown client %s", name)))
	}
	return client
}
