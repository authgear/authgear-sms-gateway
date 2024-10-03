package smsclient

import (
	"encoding/json"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sensitive"
)

type TemplateVariables struct {
	AppName     string `json:"app_name"`
	ClientID    string `json:"client_id"`
	Code        string `json:"code"`
	Email       string `json:"email"`
	HasPassword bool   `json:"has_password"`
	Host        string `json:"host"`
	Link        string `json:"link"`
	Password    string `json:"password"`
	Phone       string `json:"phone"`
	State       string `json:"state"`
	UILocales   string `json:"ui_locales"`
	URL         string `json:"url"`
	XState      string `json:"x_state"`
}

var TemplateVariablesSchema = `{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"app_name": { "type": "string" },
		"client_id": { "type": "string" },
		"code": { "type": "string" },
		"email": { "type": "string" },
		"has_password": { "type": "boolean" },
		"host": { "type": "string" },
		"link": { "type": "string" },
		"password": { "type": "string" },
		"phone": { "type": "string" },
		"state": { "type": "string" },
		"ui_locales": { "type": "string" },
		"url": { "type": "string" },
		"x_state": { "type": "string" }
	},
	"required": []
}`

type SendOptions struct {
	To                sensitive.PhoneNumber
	Body              string
	TemplateName      string
	LanguageTag       string
	TemplateVariables *TemplateVariables
}

type SendResultInfoVariable struct {
	Key         string `json:"key,omitempty"`
	ValueLength int    `json:"value_length"`
}

type SendResultInfoSendCloud struct {
	TemplateID                 string                    `json:"template_id,omitempty"`
	SendResultInfoVariableList []*SendResultInfoVariable `json:"variables,omitempty"`
}

type SendResultInfoTwilio struct {
	BodyLength   int  `json:"body_length,omitempty"`
	SegmentCount *int `json:"segment_count,omitempty"`
}

type SendResultInfoAccessYou struct{}

type SendResultInfoRoot struct {
	ProviderName string `json:"provider_name,omitempty"`
}

type SendResultInfo struct {
	SendResultInfoRoot      *SendResultInfoRoot      `json:"root,omitempty"`
	SendResultInfoTwilio    *SendResultInfoTwilio    `json:"twilio,omitempty"`
	SendResultInfoAccessYou *SendResultInfoAccessYou `json:"accessyou,omitempty"`
	SendResultInfoSendCloud *SendResultInfoSendCloud `json:"sendcloud,omitempty"`
}

type SendResult struct {
	DumpedResponse []byte          `json:"dumped_response,omitempty"`
	Success        bool            `json:"success"`
	Info           *SendResultInfo `json:"info,omitempty"`
}

func (r *SendResult) Error() string {
	jsonData, _ := json.Marshal(r)
	return string(jsonData)
}

type RawClient interface {
	Send(options *SendOptions) (*SendResult, error)
}
