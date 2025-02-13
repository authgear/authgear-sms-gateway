package handler

import (
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sensitive"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/api"
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
	// error_detail is additional information you want to let the user know
	// It may appear on the ui
	ErrorDetail    string                 `json:"error_detail,omitempty"`
	DumpedResponse []byte                 `json:"dumped_response,omitempty"`
	Info           *smsclient.SendContext `json:"info,omitempty"`
}
