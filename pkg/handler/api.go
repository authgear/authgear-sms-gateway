package handler

import (
	"net/http"

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

type Code string

const (
	// CodeOK means no error.
	CodeOK Code = "ok"

	// CodeInvalidRequest means the request is invalid.
	CodeInvalidRequest Code = "invalid_request"

	// CodeUnknownResponse means the response from the SMS provider is unknown.
	CodeUnknownResponse Code = "unknown_response"

	// CodeUnknownError means any other error.
	CodeUnknownError Code = "unknown_error"
)

func (c Code) HTTPStatusCode() int {
	switch c {
	case CodeOK:
		return http.StatusOK
	case CodeInvalidRequest:
		return http.StatusBadRequest
	case CodeUnknownError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

type ResponseBody struct {
	Code             Code                      `json:"code"`
	ErrorDescription string                    `json:"error_description,omitempty"`
	DumpedResponse   []byte                    `json:"dumped_response,omitempty"`
	Info             *smsclient.SendResultInfo `json:"info,omitempty"`
}
