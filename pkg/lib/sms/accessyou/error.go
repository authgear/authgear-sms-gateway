package accessyou

import (
	"github.com/authgear/authgear-sms-gateway/pkg/lib/api"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

func MakeError(
	msgStatus string,
	dumpedResponse []byte,
	providerName string,
) *smsclient.SendResultError {
	err := &smsclient.SendResultError{
		DumpedResponse:    dumpedResponse,
		ProviderName:      providerName,
		ProviderErrorCode: msgStatus,
	}

	// See https://www.accessyou.com/smsapi.pdf
	switch msgStatus {
	case "108":
		fallthrough
	case "110":
		err.Code = api.CodeInvalidPhoneNumber
	case "105":
		err.Code = api.CodeAuthenticationFailed
	case "106":
		fallthrough
	case "107":
		err.Code = api.CodeDeliveryRejected
	}

	return err
}
