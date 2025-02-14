package handler

import (
	"github.com/authgear/authgear-sms-gateway/pkg/lib/api"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sensitive"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type RequestBody struct {
	AppID             string                       `json:"app_id,omitempty"`
	To                sensitive.PhoneNumber        `json:"to,omitempty"`
	Body              string                       `json:"body,omitempty"`
	TemplateName      string                       `json:"template_name"`
	LanguageTag       string                       `json:"language_tag"`
	TemplateVariables *smsclient.TemplateVariables `json:"template_variables"`
}

type ResponseBody struct {
	Code             api.Code `json:"code"`
	ErrorDescription string   `json:"error_description,omitempty"`

	ProviderName      string                 `json:"provider_name,omitempty"`
	ProviderErrorCode string                 `json:"provider_error_code,omitempty"`
	DumpedResponse    []byte                 `json:"dumped_response,omitempty"`
	Info              *smsclient.SendContext `json:"info,omitempty"`
}
