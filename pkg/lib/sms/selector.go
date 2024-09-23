package sms

import (
	"errors"
	"fmt"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms"
)

type Selector struct {
	Matcher Matcher
	Client  sms.RawClient
}

type SMSProviderSelector struct {
	Selectors     []*Selector
	DefaultClient sms.RawClient
}

func NewSMSProviderSelector(c *config.SMSProviderConfig, clients *SMSProviders) (*SMSProviderSelector, error) {
	var selectors []*Selector
	var defaultClient sms.RawClient
	for _, providerSelector := range c.ProviderSelector.Switch {
		client, err := clients.GetClientByName(providerSelector.UseProvider)
		if err != nil {
			return nil, err
		}
		matcher := ParseMatcher(providerSelector)
		switch m := matcher.(type) {
		case *MatcherDefault:
			defaultClient = client
			break
		default:
			selectors = append(selectors, &Selector{
				Matcher: m,
				Client:  client,
			})
		}
	}
	return &SMSProviderSelector{
		Selectors:     selectors,
		DefaultClient: defaultClient,
	}, nil
}

func (s *SMSProviderSelector) GetClientByMatch(ctx *MatchContext) (sms.RawClient, error) {
	for _, selector := range s.Selectors {
		if selector.Matcher.Match(ctx) {
			return selector.Client, nil
		}
	}
	if s.DefaultClient != nil {
		return s.DefaultClient, nil
	}
	return nil, errors.New(fmt.Sprintf("Cannot select provider given %v", ctx))
}
