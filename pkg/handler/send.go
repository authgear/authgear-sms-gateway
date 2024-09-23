package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/authgear/authgear-server/pkg/util/httputil"

	"github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/util/validation"

	sms_infra "github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/type_util"
)

type JSONResponseWriter interface {
	WriteResponse(rw http.ResponseWriter, resp *api.Response)
}

type SendHandler struct {
	Logger     *slog.Logger
	SMSService *sms.SMSService
	JSON       JSONResponseWriter
}

func NewSendHandler(
	logger *slog.Logger,
	smsService *sms.SMSService,
	jsonResponseWriter JSONResponseWriter,
) *SendHandler {
	return &SendHandler{
		Logger:     logger,
		SMSService: smsService,
		JSON:       jsonResponseWriter,
	}
}

type RequestBody struct {
	AppID             string                         `json:"app_id,omitempty"`
	To                type_util.SensitivePhoneNumber `json:"to,omitempty"`
	Body              string                         `json:"body,omitempty"`
	TemplateName      string                         `json:"template_name"`
	LanguageTag       string                         `json:"language_tag"`
	TemplateVariables *sms_infra.TemplateVariables   `json:"template_variables"`
}

var RequestSchema = validation.NewMultipartSchema("SendRequestSchema")

var _ = RequestSchema.Add("SendRequestSchema", `
{
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"app_id": { "type": "string" },
		"to": { "type": "string" },
		"body": { "type": "string" },
		"app_id": { "type": "string" },
		"message_type": { "type": "string" },
		"template_name": { "type": "string" },
		"language_tag": { "type": "string" },
		"template_variables": { "$refs": "#/$defs/TemplateVariables" }
	},
	"required": ["to", "body", "template_name", "language_tag", "template_variables"]
}
`)
var _ = RequestSchema.Add("TemplateVariables", sms_infra.TemplateVariablesSchema)

type Result struct {
	ClientResponse string `json:"client_response"`
	Success        bool   `json:"success"`
	HasNumSegments bool   `json:"has_num_segments"`
	NumSegments    int    `json:"num_segments"`
}

func init() {
	RequestSchema.Instantiate()
}

func (h *SendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	err := httputil.BindJSONBody(r, w, RequestSchema.Validator(), &body)
	if err != nil {
		h.JSON.WriteResponse(w, &api.Response{Error: err})
		return
	}
	h.Logger.Info(fmt.Sprintf("Attempt to send sms to %v. Body: %v. AppID: %v", body.To, body.Body, body.AppID))
	sendResult, err := h.SMSService.Send(
		body.AppID,
		&sms_infra.SendOptions{
			To:                body.To,
			Body:              body.Body,
			TemplateName:      body.TemplateName,
			LanguageTag:       body.LanguageTag,
			TemplateVariables: body.TemplateVariables,
		},
	)
	if err != nil {
		h.JSON.WriteResponse(w, &api.Response{Error: err})
		return
	}

	h.JSON.WriteResponse(w, &api.Response{
		Result: Result{
			ClientResponse: string(sendResult.ClientResponse),
			Success:        sendResult.Success,
			HasNumSegments: sendResult.HasNumSegments,
			NumSegments:    sendResult.NumSegments,
		},
	})
}
