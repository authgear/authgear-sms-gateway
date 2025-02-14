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
	// CodeDeliveryRejected means the sms gateway rejected the request for some reason
	// e.g. The account was suspended
	CodeDeliveryRejected Code = "delivery_rejected"

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
	case CodeDeliveryRejected:
		return http.StatusInternalServerError
	case CodeInvalidRequest:
		return http.StatusBadRequest
	case CodeUnknownError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
