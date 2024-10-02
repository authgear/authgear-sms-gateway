package sendcloud

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

func TestMakeEffectiveTemplateVariables(t *testing.T) {
	Convey("MakeEffectiveTemplateVariables", t, func() {
		So(MakeEffectiveTemplateVariables(&smsclient.TemplateVariables{}, []*config.SendCloudTemplateVariableKeyMapping{
			{
				From: config.SendCloudTemplateVariableKeyMappingFromAppName,
				To:   "app",
			},
			{
				From: config.SendCloudTemplateVariableKeyMappingFromCode,
				To:   "mycode",
			},
		}), ShouldEqual, EffectiveTemplateVariables(map[string]interface{}{
			"app":    "",
			"mycode": "",
		}))

		So(MakeEffectiveTemplateVariables(&smsclient.TemplateVariables{
			AppName: "my-app-name",
			Code:    "123456",
		}, []*config.SendCloudTemplateVariableKeyMapping{
			{
				From: config.SendCloudTemplateVariableKeyMappingFromAppName,
				To:   "app",
			},
			{
				From: config.SendCloudTemplateVariableKeyMappingFromCode,
				To:   "somecode",
			},
		}), ShouldEqual, EffectiveTemplateVariables(map[string]interface{}{
			"app":      "my-app-name",
			"somecode": "123456",
		}))
	})
}
