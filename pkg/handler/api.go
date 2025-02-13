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

	// CodeInvalidPhoneNumber means the phone number is not a valid number
	CodeInvalidPhoneNumber Code = "invalid_phone_number"
	// CodeRateLimited means some rate limit is hit and the request should be retried later
	CodeRateLimited Code = "rate_limited"
	// CodeAuthenticationFailed means authentication failed in the sms gateway
	CodeAuthenticationFailed Code = "authentication_failed"
	// CodeAuthorizationFailed means authorization failed in the sms gateway
	CodeAuthorizationFailed Code = "authorization_failed"

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
	case CodeInvalidPhoneNumber:
		return http.StatusBadRequest
	case CodeRateLimited:
		return http.StatusTooManyRequests
	case CodeAuthenticationFailed:
		return http.StatusInternalServerError
	case CodeAuthorizationFailed:
		return http.StatusInternalServerError
	case CodeInvalidRequest:
		return http.StatusBadRequest
	case CodeUnknownError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

type ResponseBody struct {
	Code             Code   `json:"code"`
	ErrorDescription string `json:"error_description,omitempty"`
	// error_detail is additional information you want to let the user know
	// It may appear on the ui
	ErrorDetail    string                 `json:"error_detail,omitempty"`
	DumpedResponse []byte                 `json:"dumped_response,omitempty"`
	Info           *smsclient.SendContext `json:"info,omitempty"`
}
