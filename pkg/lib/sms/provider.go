package sms

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type SMSClientMap map[string]smsclient.RawClient

func NewSMSClientMap(c *config.RootConfig, httpClient *http.Client, logger *slog.Logger) SMSClientMap {
	var clientMap = make(map[string]smsclient.RawClient)

	for _, provider := range c.Providers {
		client := NewClientFromConfigProvider(provider, httpClient, logger)
		clientMap[provider.Name] = client
	}

	return SMSClientMap(clientMap)
}

func (s SMSClientMap) GetClientByName(name string) smsclient.RawClient {
	client := s[name]
	if client == nil {
		panic(fmt.Errorf("Unknown client %v", name))
	}
	return client
}
