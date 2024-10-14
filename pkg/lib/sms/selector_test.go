package sms

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
)

var config1 = `
providers:
  - name: p1
    type: twilio
    twilio:
      account_sid: "my-account-sid"
      auth_token: "my-auth-token"
      messaging_service_sid: "my-messaging-service-sid"
  - name: p2
    type: twilio
    twilio:
      account_sid: "my-account-sid"
      auth_token: "my-auth-token"
      messaging_service_sid: "my-messaging-service-sid"
  - name: p3
    type: twilio
    twilio:
      account_sid: "my-account-sid"
      auth_token: "my-auth-token"
      messaging_service_sid: "my-messaging-service-sid"
  - name: p4
    type: twilio
    twilio:
      account_sid: "my-account-sid"
      auth_token: "my-auth-token"
      messaging_service_sid: "my-messaging-service-sid"
  - name: p5
    type: twilio
    twilio:
      account_sid: "my-account-sid"
      auth_token: "my-auth-token"
      messaging_service_sid: "my-messaging-service-sid"
provider_selector:
  switch:
    - type: match_app_id_and_phone_number_alpha2
      use_provider: p1
      phone_number_alpha2: HK
      app_id: "123"
    - type: match_app_id_and_phone_number_alpha2
      use_provider: p2
      phone_number_alpha2: CN
      app_id: "123"
    - type: match_app_id
      use_provider: p3
      app_id: "use_p3"
    - type: match_phone_number_alpha2
      use_provider: p4
      phone_number_alpha2: CN
    - type: default
      use_provider: p5

`

func TestSelector(t *testing.T) {
	test := func(convey string, configYAML string, ctx *MatchContext, expectedName string) {
		c, _ := config.ParseRootConfigFromYAML([]byte(configYAML))
		res := GetClientNameByMatch(c, ctx)
		Convey(convey, t, func() {
			So(res, ShouldEqual, expectedName)
		})
	}

	test(
		"App ID and Country Code (HK) match. Pick p1",
		config1,
		&MatchContext{AppID: "123", PhoneNumber: "+85298765432"},
		"p1",
	)
	test(
		"App ID and Country Code (CN) match. Pick p2",
		config1,
		&MatchContext{AppID: "123", PhoneNumber: "+8698765432"},
		"p2",
	)
	test(
		"App ID match. Pick p3",
		config1,
		&MatchContext{AppID: "use_p3", PhoneNumber: "+8698765432"},
		"p3",
	)
	test(
		"App ID and Country Code (HK) not match. Pick default",
		config1,
		&MatchContext{AppID: "456", PhoneNumber: "+85298765432"},
		"p5",
	)
	test(
		"App ID not match. Country Code (CN) match. Pick p4",
		config1,
		&MatchContext{AppID: "456", PhoneNumber: "+8698765432"},
		"p4",
	)
}
