package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/authgear/authgear-server/pkg/util/httputil"

	"github.com/authgear/authgear-server/pkg/util/validation"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/type_util"
)

type SendHandler struct {
	Logger     *slog.Logger
	SMSService *sms.SMSService
}

func NewSendHandler(logger *slog.Logger, smsService *sms.SMSService) *SendHandler {
	return &SendHandler{
		Logger:     logger,
		SMSService: smsService,
	}
}

type RequestBody struct {
	AppID string                         `json:"app_id,omitempty"`
	To    type_util.SensitivePhoneNumber `json:"to,omitempty"`
	Body  string                         `json:"body,omitempty"`
}

var RequestSchema = validation.NewSimpleSchema(`
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"app_id": { "type": "string" },
		"to": { "type": "string" },
		"body": { "type": "string" }
	},
	"required": ["to", "body"]
}
`)

func (h *SendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	err := httputil.BindJSONBody(r, w, RequestSchema.Validator(), &body)
	if err != nil {
		panic(err)
	}
	h.Logger.Info(fmt.Sprintf("Attempt to send sms to %v. Body: %v. AppID: %v", body.To, body.Body, body.AppID))
	err = h.SMSService.Send(body.To, body.Body)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "OK")
}
