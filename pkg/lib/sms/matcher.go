package sms

import (
	"errors"
	"fmt"

	"github.com/nyaruka/phonenumbers"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
)

type MatchContext struct {
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

type MatcherDefault struct{}

func (m *MatcherDefault) Match(ctx *MatchContext) bool {
	return true
}

var _ Matcher = &MatcherDefault{}

func ParseMatcher(rule *config.ProviderSelectorSwitchRule) (Matcher, error) {
	switch rule.Type {
	case config.ProviderSelectorSwitchTypeMatchPhoneNumberAlpha2:
		return &MatcherPhoneNumberAlpha2{
			Code: rule.PhoneNumberAlpha2,
		}, nil
	case config.ProviderSelectorSwitchTypeDefault:
		return &MatcherDefault{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown rule type %s", rule.Type))
	}
}