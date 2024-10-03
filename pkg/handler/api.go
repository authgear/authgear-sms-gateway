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

type SendResultInfoAccessYou struct {
}

type SendResultInfoRoot struct {
	ProviderName string `json:"provider_name"`
}

type SendResultInfo struct {
	SendResultInfoRoot      *SendResultInfoRoot      `json:"root,omitempty"`
	SendResultInfoTwilio    *SendResultInfoTwilio    `json:"twilio,omitempty"`
	SendResultInfoAccessYou *SendResultInfoAccessYou `json:"accessyou,omitempty"`
	SendResultInfoSendCloud *SendResultInfoSendCloud `json:"sendcloud,omitempty"`
}

func MakeSendResultInfo(info *smsclient.SendResultInfo) *SendResultInfo {
	var root *SendResultInfoRoot
	var twilio *SendResultInfoTwilio
	var accessyou *SendResultInfoAccessYou
	var sendcloud *SendResultInfoSendCloud

	if info.SendResultInfoRoot != nil {
		root = (*SendResultInfoRoot)(info.SendResultInfoRoot)
	}

	if info.SendResultInfoTwilio != nil {
		twilio = (*SendResultInfoTwilio)(info.SendResultInfoTwilio)
	}

	if info.SendResultInfoAccessYou != nil {
		accessyou = (*SendResultInfoAccessYou)(info.SendResultInfoAccessYou)
	}

	if info.SendResultInfoSendCloud != nil {
		var variableList []*SendResultInfoVariable
		for _, v := range info.SendResultInfoSendCloud.SendResultInfoVariableList {
			variableList = append(variableList, (*SendResultInfoVariable)(v))
		}
		sendcloud = &SendResultInfoSendCloud{
			TemplateID:                 info.SendResultInfoSendCloud.TemplateID,
			SendResultInfoVariableList: variableList,
		}
	}

	return &SendResultInfo{
		SendResultInfoRoot:      root,
		SendResultInfoTwilio:    twilio,
		SendResultInfoAccessYou: accessyou,
		SendResultInfoSendCloud: sendcloud,
	}
}

type ResponseBody struct {
	Code             Code   `json:"code"`
	ErrorDescription string `json:"error_description,omitempty"`
	DumpedResponse   []byte `json:"dumped_response,omitempty"`

	Info *SendResultInfo `json:"info,omitempty"`
}
