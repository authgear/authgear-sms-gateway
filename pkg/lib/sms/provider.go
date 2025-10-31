package sms

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type SMSProviderMap map[string]smsclient.RawClient

func NewSMSProviderMap(c *config.RootConfig, httpClient *http.Client, logger *slog.Logger) SMSProviderMap {
	var clientMap = make(map[string]smsclient.RawClient)

	for _, provider := range c.Providers {
		client := NewClientFromConfigProvider(provider, httpClient, logger)
		clientMap[provider.Name] = client
	}

	return SMSProviderMap(clientMap)
}

func (s SMSProviderMap) GetProviderByName(name string) smsclient.RawClient {
	client := s[name]
	if client == nil {
		panic(fmt.Errorf("unknown client %v", name))
	}
	return client
}
