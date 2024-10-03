package sendcloud

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

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

func (n *SendCloudClient) Send(options *smsclient.SendOptions) (*smsclient.SendResult, error) {
	info := &smsclient.SendResultInfo{
		SendResultInfoSendCloud: &smsclient.SendResultInfoSendCloud{},
	}

	template, err := n.TemplateResolver.Resolve(options.TemplateName, options.LanguageTag)
	if err != nil {
		return nil, err
	}
	info.SendResultInfoSendCloud.TemplateID = string(template.TemplateID)
	templateVariables := MakeEffectiveTemplateVariables(options.TemplateVariables, template.TemplateVariableKeyMappings)

	var sendResultInfoVariableList []*smsclient.SendResultInfoVariable

	for key, value := range templateVariables {
		sendResultInfoVariableList = append(sendResultInfoVariableList, &smsclient.SendResultInfoVariable{
			Key:         key,
			ValueLength: len(fmt.Sprintf("%v", value)),
		})
	}
	info.SendResultInfoSendCloud.SendResultInfoVariableList = sendResultInfoVariableList

	sendRequest := NewSendRequest(
		string(template.TemplateMsgType),
		[]string{
			string(options.To),
		},
		n.SMSUser,
		string(template.TemplateID),
		templateVariables.WrapKeys(),
	)

	dumpedResponse, sendResponse, err := Send(n.Client, n.BaseUrl, &sendRequest, n.SMSKey, n.Logger)
	if err != nil {
		return nil, err
	}

	return &smsclient.SendResult{
		DumpedResponse: dumpedResponse,
		Success:        sendResponse.StatusCode == 200,
		Info:           info,
	}, nil
}

var _ smsclient.RawClient = &SendCloudClient{}
