package config

type SendCloudTemplateVariableKeyMappingFrom string

const (
	SendCloudTemplateVariableKeyMappingFromAppName     SendCloudTemplateVariableKeyMappingFrom = "app_name"
	SendCloudTemplateVariableKeyMappingFromClientID    SendCloudTemplateVariableKeyMappingFrom = "client_id"
	SendCloudTemplateVariableKeyMappingFromCode        SendCloudTemplateVariableKeyMappingFrom = "code"
	SendCloudTemplateVariableKeyMappingFromEmail       SendCloudTemplateVariableKeyMappingFrom = "email"
	SendCloudTemplateVariableKeyMappingFromHasPassword SendCloudTemplateVariableKeyMappingFrom = "has_password"
	SendCloudTemplateVariableKeyMappingFromHost        SendCloudTemplateVariableKeyMappingFrom = "host"
	SendCloudTemplateVariableKeyMappingFromLink        SendCloudTemplateVariableKeyMappingFrom = "link"
	SendCloudTemplateVariableKeyMappingFromPassword    SendCloudTemplateVariableKeyMappingFrom = "password"
	SendCloudTemplateVariableKeyMappingFromPhone       SendCloudTemplateVariableKeyMappingFrom = "phone"
	SendCloudTemplateVariableKeyMappingFromState       SendCloudTemplateVariableKeyMappingFrom = "state"
	SendCloudTemplateVariableKeyMappingFromUILocales   SendCloudTemplateVariableKeyMappingFrom = "ui_locales"
	SendCloudTemplateVariableKeyMappingFromURL         SendCloudTemplateVariableKeyMappingFrom = "url"
	SendCloudTemplateVariableKeyMappingFromXState      SendCloudTemplateVariableKeyMappingFrom = "x_state"
)

var _ = RootSchema.Add("SendCloudTemplateVariableKeyMappingFrom", `
{
	"type": "string",
	"enum": [
		"app_name",
		"client_id",
		"code",
		"email",
		"has_password",
		"host",
		"link",
		"password",
		"phone",
		"state",
		"ui_locales",
		"url",
		"x_state"
	]
}`)

type SendCloudTemplateVariableKeyMapping struct {
	From SendCloudTemplateVariableKeyMappingFrom `json:"from,omitempty"`
	To   string                                  `json:"to,omitempty"`
}

var _ = RootSchema.Add("SendCloudTemplateVariableKeyMapping", `
{
  "type": "object",
	"additionalProperties": false,
	"properties": {
		"from": { "$ref": "#/$defs/SendCloudTemplateVariableKeyMappingFrom" },
		"to": { "type": "string" }
	},
	"required": ["from", "to"]
}`)

type SendCloudTemplate struct {
	TemplateID                  string                                 `json:"template_id,omitempty"`
	TemplateMsgType             string                                 `json:"template_msg_type,omitempty"`
	TemplateVariableKeyMappings []*SendCloudTemplateVariableKeyMapping `json:"template_variable_key_mappings,omitempty"`
}

var _ = RootSchema.Add("SendCloudTemplate", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"template_id": { "type": "string" },
		"template_msg_type": { "type": "string" },
		"template_variable_key_mappings": {
			"type": "array",
			"minItems": 1,
			"items": { "$ref": "#/$defs/SendCloudTemplateVariableKeyMapping" }
		}
	},
	"required": ["template_id", "template_msg_type", "template_variable_key_mappings"]
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
	"required": ["sms_user", "sms_key", "templates", "template_assignments"]
}
`)
