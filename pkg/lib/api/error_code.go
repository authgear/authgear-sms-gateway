package api

import "net/http"

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
	// CodeAttemptedToSendOTPTemplateWithoutCode means the message is not otp message, and the provider only support otp message
	CodeAttemptedToSendOTPTemplateWithoutCode Code = "attempted_to_send_otp_template_without_code"
	// CodeDeliveryRejected means the sms gateway rejected the request for some reason
	// e.g. The account was suspended
	CodeDeliveryRejected Code = "delivery_rejected"

	// CodeInvalidRequest means the request is invalid.
	CodeInvalidRequest Code = "invalid_request"

	// CodeUnknownResponse means the response from the SMS provider is unknown.
	CodeUnknownResponse Code = "unknown_response"

	// CodeUnknownError means any other error.
	CodeUnknownError Code = "unknown_error"

	// Error codes we use internally
	// CodeUnknownResponseFormat means the response from the SMS provider is unknown format.
	CodeUnknownResponseFormat Code = "_unknown_response_format"
	// CodeProviderTimeout means the SMS provider timed out.
	CodeProviderTimeout Code = "_provider_timeout"
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
	case CodeAttemptedToSendOTPTemplateWithoutCode:
		return http.StatusBadRequest
	case CodeDeliveryRejected:
		return http.StatusInternalServerError
	case CodeInvalidRequest:
		return http.StatusBadRequest
	case CodeUnknownError:
		return http.StatusInternalServerError
	// Internal codes
	case CodeUnknownResponseFormat:
		return http.StatusInternalServerError
	case CodeProviderTimeout:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
