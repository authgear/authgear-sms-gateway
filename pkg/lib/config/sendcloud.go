package config

type SendCloudTemplate struct {
	TemplateID      string `json:"template_id,omitempty"`
	TemplateMsgType string `json:"template_msg_type,omitempty"`
}

var _ = RootSchema.Add("SendCloudTemplate", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"template_id": { "type": "string" },
		"template_msg_type": { "type": "string" }
	},
	"required": ["template_id", "template_msg_type"]
}
`)

type SendCloudTemplateAssignmentByLanguage struct {
	AuthgearLanguage string `json:"authgear_language,omitempty"`
	TemplateID       string `json:"template_id,omitempty"`
}

var _ = RootSchema.Add("SendCloudTemplateAssignmentByLanguage", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"authgear_language": { "type": "string" },
		"template_id": { "type": "string" }
	},
	"required": ["authgear_language", "template_id"]
}
`)

type SendCloudTemplateAssignment struct {
	AuthgearTemplateName string                                   `json:"authgear_template_name,omitempty"`
	DefaultTemplateID    string                                   `json:"default_template_id,omitempty"`
	ByLanguages          []*SendCloudTemplateAssignmentByLanguage `json:"by_languages,omitempty"`
}

var _ = RootSchema.Add("SendCloudTemplateAssignment", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"authgear_template_name": { "type": "string" },
		"default_template_id": { "type": "string" },
		"by_languages": {
			"type": "array",
			"minItems": 1,
			"items": { "$ref": "#/defs/SendCloudTemplateAssignmentByLanguage"}
		}
	},
	"required": ["authgear_template_name", "default_template_idd"]
}
`)

type ProviderConfigSendCloud struct {
	Sender              string                         `json:"sender,omitempty"`
	BaseUrl             string                         `json:"base_url,omitempty"`
	SMSUser             string                         `json:"sms_user,omitempty"`
	SMSKey              string                         `json:"sms_key,omitempty"`
	Templates           []*SendCloudTemplate           `json:"templates,omitempty"`
	TemplateAssignments []*SendCloudTemplateAssignment `json:"template_assignments,omitempty"`
}

var _ = RootSchema.Add("ProviderConfigSendCloud", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"sender": { "type": "string" },
		"base_url": { "type": "string" },
		"sms_user": { "type": "string" },
		"sms_key": {"type": "string"},
		"templates": {
			"type": "array",
			"minItems": 1,
			"items": { "$refs": "#/$defs/SendCloudTemplate" }
		},
		"template_assignments": {
			"type": "array",
			"minItems": 1,
			"items": { "$refs": "#/$defs/SendCloudTemplateAssignment" }
		}
	},
	"required": ["sender", "sms_user", "sms_key", "templates", "template_assignments"]
}
`)
