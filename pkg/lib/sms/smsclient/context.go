package smsclient

import (
	"context"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sensitive"
)

type SendContext struct {
	Root      *SendContextRoot      `json:"root,omitempty"`
	Twilio    *SendContextTwilio    `json:"twilio,omitempty"`
	AccessYou *SendContextAccessYou `json:"accessyou,omitempty"`
	SendCloud *SendContextSendCloud `json:"sendcloud,omitempty"`
}

type SendContextRoot struct {
	AppID        string                `json:"app_id,omitempty"`
	To           sensitive.PhoneNumber `json:"to,omitempty"`
	TemplateName string                `json:"template_name,omitempty"`
	LanguageTag  string                `json:"language_tag,omitempty"`

	ProviderName string `json:"provider_name,omitempty"`
}

type SendContextAccessYou struct{}

type SendContextTwilio struct {
	BodyLength   int  `json:"body_length,omitempty"`
	SegmentCount *int `json:"segment_count,omitempty"`
}

type SendContextSendCloud struct {
	TemplateID string                 `json:"template_id,omitempty"`
	Variables  []*SendContextVariable `json:"variables,omitempty"`
}

type SendContextVariable struct {
	Key         string `json:"key,omitempty"`
	ValueLength int    `json:"value_length"`
}

type sendContextKeyType struct{}

var sendContextKey = sendContextKeyType{}

func WithSendContext(ctx context.Context, f func(sendCtx *SendContext)) context.Context {
	sendContext, ok := ctx.Value(sendContextKey).(*SendContext)
	if !ok || sendContext == nil {
		sendContext = &SendContext{}
	}
	f(sendContext)
	return context.WithValue(ctx, sendContextKey, sendContext)
}

func GetSendContext(ctx context.Context) *SendContext {
	sendContext, ok := ctx.Value(sendContextKey).(*SendContext)
	if !ok || sendContext == nil {
		return &SendContext{}
	}

	return sendContext
}
