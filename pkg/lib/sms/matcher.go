package sms

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
)

type MatchContext struct {
	AppID       string
	PhoneNumber string
}

type Matcher interface {
	Match(ctx *MatchContext) bool
}

type MatcherPhoneNumberAlpha2 struct {
	Code string
}

var _ Matcher = &MatcherPhoneNumberAlpha2{}

func (m *MatcherPhoneNumberAlpha2) Match(ctx *MatchContext) bool {
	num, err := phonenumbers.Parse(ctx.PhoneNumber, "")
	if err != nil {
		return false
	}
	regionCode := phonenumbers.GetRegionCodeForNumber(num)
	return m.Code == regionCode
}

type MatcherAppID struct {
	AppID string
}

var _ Matcher = &MatcherAppID{}

func (m *MatcherAppID) Match(ctx *MatchContext) bool {
	if m.AppID == "" {
		return false
	}

	if m.AppID == ctx.AppID {
		return true
	}

	return false
}

type MatcherAppIDAndPhoneNumberAlpha2 struct {
	AppID string
	Code  string
}

var _ Matcher = &MatcherAppIDAndPhoneNumberAlpha2{}

func (m *MatcherAppIDAndPhoneNumberAlpha2) Match(ctx *MatchContext) bool {
	num, err := phonenumbers.Parse(ctx.PhoneNumber, "")
	if err != nil {
		return false
	}
	regionCode := phonenumbers.GetRegionCodeForNumber(num)
	if m.AppID == "" {
		// Any app id if app id is not specified in config
		return m.Code == regionCode
	}
	return m.Code == regionCode && m.AppID == ctx.AppID
}

type MatcherDefault struct{}

var _ Matcher = &MatcherDefault{}

func (m *MatcherDefault) Match(ctx *MatchContext) bool {
	return true
}

func ParseMatcher(rule *config.ProviderSelectorSwitchRule) Matcher {
	switch rule.Type {
	case config.ProviderSelectorSwitchTypeMatchPhoneNumberAlpha2:
		return &MatcherPhoneNumberAlpha2{
			Code: rule.PhoneNumberAlpha2,
		}
	case config.ProviderSelectorSwitchTypeMatchAppID:
		return &MatcherAppID{
			AppID: rule.AppID,
		}
	case config.ProviderSelectorSwitchTypeMatchAppIDAndPhoneNumberAlpha2:
		return &MatcherAppIDAndPhoneNumberAlpha2{
			AppID: rule.AppID,
			Code:  rule.PhoneNumberAlpha2,
		}
	case config.ProviderSelectorSwitchTypeDefault:
		return &MatcherDefault{}
	default:
		panic(fmt.Errorf("unknown rule type %s", rule.Type))
	}
}
