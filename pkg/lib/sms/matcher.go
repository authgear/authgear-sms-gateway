package sms

import (
	"errors"
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

func (m *MatcherPhoneNumberAlpha2) Match(ctx *MatchContext) bool {
	num, err := phonenumbers.Parse(ctx.PhoneNumber, "")
	if err != nil {
		return false
	}
	regionCode := phonenumbers.GetRegionCodeForNumber(num)
	return m.Code == regionCode
}

var _ Matcher = &MatcherPhoneNumberAlpha2{}

type MatcherAppIDAndPhoneNumberAlpha2 struct {
	AppID string
	Code  string
}

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

var _ Matcher = &MatcherAppIDAndPhoneNumberAlpha2{}

type MatcherDefault struct{}

func (m *MatcherDefault) Match(ctx *MatchContext) bool {
	return true
}

var _ Matcher = &MatcherDefault{}

func ParseMatcher(rule *config.ProviderSelectorSwitchRule) Matcher {
	switch rule.Type {
	case config.ProviderSelectorSwitchTypeMatchPhoneNumberAlpha2:
		return &MatcherPhoneNumberAlpha2{
			Code: rule.PhoneNumberAlpha2,
		}
	case config.ProviderSelectorSwitchTypeMatchAppIDAndPhoneNumberAlpha2:
		return &MatcherAppIDAndPhoneNumberAlpha2{
			AppID: rule.AppID,
			Code:  rule.PhoneNumberAlpha2,
		}
	case config.ProviderSelectorSwitchTypeDefault:
		return &MatcherDefault{}
	default:
		panic(errors.New(fmt.Sprintf("Unknown rule type %s", rule.Type)))
	}
}
