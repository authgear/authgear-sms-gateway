package api

import "net/http"

type Code string

// Read authgear-server/docs/specs/sms_gateway.md
const (
	CodeOK Code = "ok"

	CodeInvalidPhoneNumber   Code = "invalid_phone_number"
	CodeRateLimited          Code = "rate_limited"
	CodeAuthenticationFailed Code = "authentication_failed"
	CodeUnsupportedRequest   Code = "unsupported_request"
	CodeDeliveryRejected     Code = "delivery_rejected"
	CodeTimeout              Code = "timeout"

	CodeInvalidRequest Code = "invalid_request"

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
	case CodeUnsupportedRequest:
		return http.StatusBadRequest
	case CodeDeliveryRejected:
		return http.StatusInternalServerError
	case CodeTimeout:
		return http.StatusInternalServerError
	case CodeInvalidRequest:
		return http.StatusBadRequest
	case CodeUnknownError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
