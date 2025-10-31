package smsclient

import (
	"context"
	"log/slog"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/logger"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sensitive"
)

type SendContext struct {
	Root      *SendContextRoot      `json:"root,omitempty"`
	Twilio    *SendContextTwilio    `json:"twilio,omitempty"`
	AccessYou *SendContextAccessYou `json:"accessyou,omitempty"`
	SendCloud *SendContextSendCloud `json:"sendcloud,omitempty"`
}

var _ logger.LoggerContexter = &SendContext{}

func (c *SendContext) GetAttrs() []slog.Attr {
	var attrs []slog.Attr

	if c.Root != nil {
		if c.Root.AppID != "" {
			attrs = append(attrs, slog.String("app_id", c.Root.AppID))
		}
		if c.Root.To != "" {
			attrs = append(attrs, slog.Any("to", c.Root.To))
		}
		if c.Root.TemplateName != "" {
			attrs = append(attrs, slog.String("template_name", c.Root.TemplateName))
		}
		if c.Root.LanguageTag != "" {
			attrs = append(attrs, slog.String("language_tag", c.Root.LanguageTag))
		}
		if c.Root.ProviderName != "" {
			attrs = append(attrs, slog.String("provider_name", c.Root.ProviderName))
		}
	}

	if c.Twilio != nil {
		if c.Twilio.BodyLength != 0 {
			attrs = append(attrs, slog.Int("body_length", c.Twilio.BodyLength))
		}
		if c.Twilio.SegmentCount != nil {
			attrs = append(attrs, slog.Int("segment_count", *c.Twilio.SegmentCount))
		}
	}

	return attrs
}

type SendContextRoot struct {
	AppID        string                `json:"app_id,omitempty"`
	To           sensitive.PhoneNumber `json:"to,omitempty"`
	TemplateName string                `json:"template_name,omitempty"`
	LanguageTag  string                `json:"language_tag,omitempty"`

	ProviderName string `json:"provider_name,omitempty"`
	ProviderType string `json:"provider_type,omitempty"`
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

type SendContextKeyType struct{}

var SendContextKey = SendContextKeyType{}

func WithSendContext(ctx context.Context, f func(sendCtx *SendContext)) context.Context {
	sendContext, ok := ctx.Value(SendContextKey).(*SendContext)
	if !ok || sendContext == nil {
		sendContext = &SendContext{}
	}
	f(sendContext)
	return context.WithValue(ctx, SendContextKey, sendContext)
}

func GetSendContext(ctx context.Context) *SendContext {
	sendContext, ok := ctx.Value(SendContextKey).(*SendContext)
	if !ok || sendContext == nil {
		return &SendContext{}
	}

	return sendContext
}
