package sendcloud

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
)

func TestTemplateResolver(t *testing.T) {
	Convey("TemplateResolver error", t, func() {
		templateResolver, err := NewSendCloudTemplateResolver(
			[]*config.SendCloudTemplate{
				&config.SendCloudTemplate{
					TemplateID:      "919880",
					TemplateMsgType: "2",
				},
				&config.SendCloudTemplate{
					TemplateID:      "919877",
					TemplateMsgType: "2",
				},
			},
			[]*config.SendCloudTemplateAssignment{
				&config.SendCloudTemplateAssignment{
					AuthgearTemplateName: "verification_sms.txt",
					DefaultTemplateID:    "919880",
					ByLanguages: []*config.SendCloudTemplateAssignmentByLanguage{
						&config.SendCloudTemplateAssignmentByLanguage{
							AuthgearLanguage: "zh",
							TemplateID:       "919878",
						},
					},
				},
			},
		)

		So(err.Error(), ShouldEqual, "Cannot find template with id 919878")
		So(templateResolver, ShouldBeNil)
	})

	Convey("TemplateResolver", t, func() {
		templateResolver, err := NewSendCloudTemplateResolver(
			[]*config.SendCloudTemplate{
				&config.SendCloudTemplate{
					TemplateID:      "919880",
					TemplateMsgType: "2",
				},
				&config.SendCloudTemplate{
					TemplateID:      "919877",
					TemplateMsgType: "2",
				},
			},
			[]*config.SendCloudTemplateAssignment{
				&config.SendCloudTemplateAssignment{
					AuthgearTemplateName: "verification_sms.txt",
					DefaultTemplateID:    "919880",
					ByLanguages: []*config.SendCloudTemplateAssignmentByLanguage{
						&config.SendCloudTemplateAssignmentByLanguage{
							AuthgearLanguage: "zh",
							TemplateID:       "919877",
						},
					},
				},
			},
		)

		So(err, ShouldBeNil)

		Convey("Should resolve", func() {
			template, err := templateResolver.Resolve("verification_sms.txt", "zh")
			So(err, ShouldBeNil)
			So(template.TemplateID, ShouldEqual, TemplateID("919877"))
			So(template.TemplateMsgType, ShouldEqual, TemplateMessageType("2"))
		})

		Convey("Should resolve default", func() {
			template, err := templateResolver.Resolve("verification_sms.txt", "gg")
			So(err, ShouldBeNil)
			So(template.TemplateID, ShouldEqual, TemplateID("919880"))
			So(template.TemplateMsgType, ShouldEqual, TemplateMessageType("2"))
		})

		Convey("Should be error due to template name not found", func() {
			template, err := templateResolver.Resolve("gg.txt", "zh")
			So(err.Error(), ShouldEqual, "Could not found template assignment from template name gg.txt")
			So(template, ShouldBeNil)
		})
	})
}
