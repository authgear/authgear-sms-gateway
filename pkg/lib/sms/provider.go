package sms

import (
	"errors"
	"fmt"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms"
)

type SMSProviders struct {
	Clients []sms.RawClient
	Map     map[string]sms.RawClient
}

func NewSMSProviders(c *config.SMSProviderConfig) (*SMSProviders, error) {
	var clients []sms.RawClient
	var clientMap = make(map[string]sms.RawClient)

	for _, provider := range c.Providers {
		client, err := sms.NewClientFromConfigProvider(provider)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
		clientMap[provider.Name] = client
	}

	return &SMSProviders{
		Clients: clients,
		Map:     clientMap,
	}, nil
}

func (s *SMSProviders) GetClientByName(name string) (sms.RawClient, error) {
	client, exists := s.Map[name]
	if !exists {
		return nil, errors.New(fmt.Sprintf("Unknown client %s", name))
	}
	return client, nil
}