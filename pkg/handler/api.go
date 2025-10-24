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
	// These will be included in apierror
	Code              api.Code `json:"code"`
	ProviderName      string   `json:"provider_name,omitempty"`
	ProviderType      string   `json:"provider_type,omitempty"`
	ProviderErrorCode string   `json:"provider_error_code,omitempty"`

	// These are only in debug logs
	GoError        string                 `json:"go_error,omitempty"`
	DumpedResponse []byte                 `json:"dumped_response,omitempty"`
	Info           *smsclient.SendContext `json:"info,omitempty"`
	// This flag is used to tell authgear server that this error is safe to ignore
	IsNonCritical bool `json:"is_non_critical,omitempty"`
}
