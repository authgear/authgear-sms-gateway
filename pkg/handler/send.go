package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/authgear/authgear-server/pkg/util/httputil"
	"github.com/authgear/authgear-server/pkg/util/validation"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type SendHandler struct {
	Logger     *slog.Logger
	SMSService *sms.SMSService
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
		"template_name": { "type": "string" },
		"language_tag": { "type": "string" },
		"template_variables": { "$refs": "#/$defs/TemplateVariables" }
	},
	"required": ["to", "body", "template_name", "language_tag", "template_variables"]
}
`)
var _ = RequestSchema.Add("TemplateVariables", smsclient.TemplateVariablesSchema)

func init() {
	RequestSchema.Instantiate()
}

func (h *SendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	err := httputil.BindJSONBody(r, w, RequestSchema.Validator(), &body)
	if err != nil {
		h.write(w, &ResponseBody{
			Code:             CodeInvalidRequest,
			ErrorDescription: err.Error(),
		})
		return
	}

	h.Logger.Info("received send request",
		"app_id", body.AppID,
		"to", body.To,
		"template_name", body.TemplateName,
		"language_tag", body.LanguageTag,
	)

	sendResult, err := h.SMSService.Send(
		body.AppID,
		&smsclient.SendOptions{
			To:                body.To,
			Body:              body.Body,
			TemplateName:      body.TemplateName,
			LanguageTag:       body.LanguageTag,
			TemplateVariables: body.TemplateVariables,
		},
	)
	// TODO: handle err.
	if err != nil {
		h.write(w, &ResponseBody{
			Code:             CodeUnknownError,
			ErrorDescription: err.Error(),
		})
		return
	}

	h.write(w, &ResponseBody{
		Code:                       CodeOK,
		UnderlyingHTTPResponseBody: string(sendResult.ClientResponse),
		SegmentCount:               sendResult.SegmentCount,
	})
}

func (h *SendHandler) write(w http.ResponseWriter, body *ResponseBody) {
	statusCode := body.Code.HTTPStatusCode()
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := encoder.Encode(body)
	if err != nil {
		panic(err)
	}
}
