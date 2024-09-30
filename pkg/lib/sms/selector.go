package sms

import (
	"fmt"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
)

func GetClientNameByMatch(c *config.RootConfig, ctx *MatchContext) string {
	var defaultClient string
	for _, providerSelector := range c.ProviderSelector.Switch {
		matcher := ParseMatcher(providerSelector)
		switch m := matcher.(type) {
		case *MatcherDefault:
			defaultClient = providerSelector.UseProvider
		default:
			if m.Match(ctx) {
				return providerSelector.UseProvider
			}
		}
	}
	if defaultClient == "" {
		panic(fmt.Errorf("Cannot select provider given %v", ctx))
	}
	return defaultClient
}
