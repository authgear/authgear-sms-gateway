package sendcloud

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/api"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

// Sendcloud will return status code 412 if the phone number is a +86 number with "+86" prefix.
// It only accepts "+" and country calling code if the phone number is NOT a +86 number.
func fixPhoneNumber(e164 string) string {
	return strings.TrimPrefix(e164, "+86")
}

type EffectiveTemplateVariables map[string]interface{}

func wrapped(field string) string {
	return fmt.Sprintf("%%%v%%", field)
}

func (v EffectiveTemplateVariables) WrapKeys() map[string]interface{} {
	res := make(map[string]interface{})
	for key, value := range v {
		res[wrapped(key)] = value
	}
	return res
}

func MakeEffectiveTemplateVariables(variables *smsclient.TemplateVariables, mappings []*config.SendCloudTemplateVariableKeyMapping) EffectiveTemplateVariables {
	res := make(EffectiveTemplateVariables)
	for _, mapping := range mappings {
		switch mapping.From {
		case config.SendCloudTemplateVariableKeyMappingFromAppName:
			res[mapping.To] = variables.AppName
		case config.SendCloudTemplateVariableKeyMappingFromClientID:
			res[mapping.To] = variables.ClientID
		case config.SendCloudTemplateVariableKeyMappingFromCode:
			res[mapping.To] = variables.Code
		case config.SendCloudTemplateVariableKeyMappingFromEmail:
			res[mapping.To] = variables.Email
		case config.SendCloudTemplateVariableKeyMappingFromHasPassword:
			res[mapping.To] = variables.HasPassword
		case config.SendCloudTemplateVariableKeyMappingFromHost:
			res[mapping.To] = variables.Host
		case config.SendCloudTemplateVariableKeyMappingFromLink:
			res[mapping.To] = variables.Link
		case config.SendCloudTemplateVariableKeyMappingFromPassword:
			res[mapping.To] = variables.Password
		case config.SendCloudTemplateVariableKeyMappingFromPhone:
			res[mapping.To] = variables.Phone
		case config.SendCloudTemplateVariableKeyMappingFromState:
			res[mapping.To] = variables.State
		case config.SendCloudTemplateVariableKeyMappingFromUILocales:
			res[mapping.To] = variables.UILocales
		case config.SendCloudTemplateVariableKeyMappingFromURL:
			res[mapping.To] = variables.URL
		case config.SendCloudTemplateVariableKeyMappingFromXState:
			res[mapping.To] = variables.XState
		}
	}
	return res
}

type SendCloudClient struct {
	BaseUrl          string
	Client           *http.Client
	SMSUser          string
	SMSKey           string
	TemplateResolver ISendCloudTemplateResolver
	Logger           *slog.Logger
}

func NewSendCloudClient(
	httpClient *http.Client,
	baseUrl string,
	smsUser string,
	smsKey string,
	templateResolver *SendCloudTemplateResolver,
	logger *slog.Logger,
) *SendCloudClient {
	if baseUrl == "" {
		baseUrl = "https://api.sendcloud.net"
	}
	return &SendCloudClient{
		Client:           httpClient,
		BaseUrl:          baseUrl,
		SMSUser:          smsUser,
		SMSKey:           smsKey,
		TemplateResolver: templateResolver,
		Logger:           logger,
	}
}

func (n *SendCloudClient) Send(ctx context.Context, options *smsclient.SendOptions) (*smsclient.SendResultSuccess, error) {
	template, err := n.TemplateResolver.Resolve(options.TemplateName, options.LanguageTag)
	if err != nil {
		return nil, err
	}

	templateVariables := MakeEffectiveTemplateVariables(options.TemplateVariables, template.TemplateVariableKeyMappings)
	var variables []*smsclient.SendContextVariable
	for key, value := range templateVariables {
		variables = append(variables, &smsclient.SendContextVariable{
			Key:         key,
			ValueLength: len(fmt.Sprintf("%v", value)),
		})
	}

	ctx = smsclient.WithSendContext(ctx, func(sendCtx *smsclient.SendContext) {
		sendCtx.SendCloud = &smsclient.SendContextSendCloud{
			TemplateID: string(template.TemplateID),
			Variables:  variables,
		}
	})

	sendRequest := NewSendRequest(
		string(template.TemplateMsgType),
		[]string{
			fixPhoneNumber(string(options.To)),
		},
		n.SMSUser,
		string(template.TemplateID),
		templateVariables.WrapKeys(),
	)

	dumpedResponse, sendResponse, err := Send(ctx, n.Client, n.BaseUrl, &sendRequest, n.SMSKey, n.Logger)
	if err != nil {
		return nil, err
	}

	// Success case.
	if sendResponse.StatusCode == 200 {
		return &smsclient.SendResultSuccess{
			DumpedResponse: dumpedResponse,
		}, nil
	} else {
		return nil, n.makeError(sendResponse.StatusCode, dumpedResponse)
	}
}

func (t *SendCloudClient) makeError(
	statusCode int,
	dumpedResponse []byte,
) *smsclient.SendResultError {
	err := &smsclient.SendResultError{
		DumpedResponse:    dumpedResponse,
		ProviderName:      "sendcloud",
		ProviderErrorCode: fmt.Sprintf("%d", statusCode),
	}

	// See https://www.sendcloud.net/doc/sms/api/
	switch statusCode {
	case 412:
		err.Code = api.CodeInvalidPhoneNumber
	case 50000:
		err.Code = api.CodeRateLimited
	case 422:
		fallthrough
	case 471:
		fallthrough
	case 474:
		err.Code = api.CodeAuthenticationFailed
	case 499:
		fallthrough
	case 473:
		err.Code = api.CodeDeliveryRejected
	}

	return err
}

var _ smsclient.RawClient = &SendCloudClient{}
