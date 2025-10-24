package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/authgear/authgear-server/pkg/util/httputil"
	"github.com/authgear/authgear-server/pkg/util/validation"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/api"
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
			Code:    api.CodeInvalidRequest,
			GoError: err.Error(),
		})
		return
	}

	r = r.WithContext(smsclient.WithSendContext(
		r.Context(),
		func(sendCtx *smsclient.SendContext) {
			sendCtx.Root = &smsclient.SendContextRoot{
				AppID:        body.AppID,
				To:           body.To,
				TemplateName: body.TemplateName,
				LanguageTag:  body.LanguageTag,
			}
		},
	))

	h.Logger.InfoContext(r.Context(), "received send request")

	sendResult, err := h.SMSService.Send(
		r.Context(),
		body.AppID,
		&smsclient.SendOptions{
			To:                body.To,
			Body:              body.Body,
			TemplateName:      body.TemplateName,
			LanguageTag:       body.LanguageTag,
			TemplateVariables: body.TemplateVariables,
		},
	)

	if err != nil {
		var errorUnsuccessResponse *smsclient.SendResultError
		if errors.As(err, &errorUnsuccessResponse) {
			h.Logger.ErrorContext(r.Context(), "unsuccessful response",
				"dumped_response", string(errorUnsuccessResponse.DumpedResponse),
				"error", err.Error(),
			)
			info := smsclient.GetSendContext(r.Context())
			code := api.CodeUnknownResponse
			if errorUnsuccessResponse.Code != "" {
				code = errorUnsuccessResponse.Code
			}
			h.write(w, &ResponseBody{
				Code:              code,
				ProviderName:      info.Root.ProviderName,
				ProviderType:      info.Root.ProviderType,
				ProviderErrorCode: errorUnsuccessResponse.ProviderErrorCode,
				DumpedResponse:    errorUnsuccessResponse.DumpedResponse,
				GoError:           err.Error(),
				Info:              info,
			})
			return
		}

		h.Logger.ErrorContext(r.Context(), "unknown error",
			"error", err.Error(),
		)
		h.write(w, &ResponseBody{
			Code:    api.CodeUnknownError,
			GoError: err.Error(),
		})
		return
	}

	h.Logger.InfoContext(r.Context(), "finished send request")

	info := smsclient.GetSendContext(r.Context())
	h.write(w, &ResponseBody{
		Code:           api.CodeOK,
		ProviderName:   info.Root.ProviderName,
		ProviderType:   info.Root.ProviderType,
		DumpedResponse: sendResult.DumpedResponse,
		Info:           info,
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
