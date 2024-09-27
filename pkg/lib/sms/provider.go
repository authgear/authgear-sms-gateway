package sms

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type SMSClientMap map[string]smsclient.RawClient

func NewSMSClientMap(c *config.RootConfig, logger *slog.Logger) SMSClientMap {
	var clientMap = make(map[string]smsclient.RawClient)

	for _, provider := range c.Providers {
		client := NewClientFromConfigProvider(provider, logger)
		clientMap[provider.Name] = client
	}

	return SMSClientMap(clientMap)
}

func (s SMSClientMap) GetClientByName(name string) smsclient.RawClient {
	client := s[name]
	if client == nil {
		panic(errors.New(fmt.Sprintf("Unknown client %s", name)))
	}
	return client
}
