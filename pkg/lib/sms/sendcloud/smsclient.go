package sendcloud

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

var ErrMissingSendCloudConfiguration = errors.New("accessyou: configuration is missing")

func makeVarsFromTemplateVariables(variables *smsclient.TemplateVariables) map[string]interface{} {
	wrapped := func(field string) string {
		return fmt.Sprintf("%%%v%%", field)
	}
	wrapKeys := func(obj map[string]interface{}) map[string]interface{} {
		res := make(map[string]interface{})
		for key, value := range obj {
			res[wrapped(key)] = value
		}
		return res
	}
	return wrapKeys(map[string]interface{}{
		"app":          variables.AppName,
		"client_id":    variables.ClientID,
		"code":         variables.Code,
		"email":        variables.Email,
		"has_password": variables.HasPassword,
		"host":         variables.Host,
		"link":         variables.Link,
		"password":     variables.Password,
		"phone":        variables.Phone,
		"state":        variables.State,
		"ui_locales":   variables.UILocales,
		"url":          variables.URL,
		"x_state":      variables.XState,
	})
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
	template, err := n.TemplateResolver.Resolve(options.TemplateName, options.LanguageTag)
	if err != nil {
		return nil, err
	}
	sendRequest := NewSendRequest(
		string(template.TemplateMsgType),
		[]string{
			string(options.To),
		},
		n.SMSUser,
		string(template.TemplateID),
		makeVarsFromTemplateVariables(options.TemplateVariables),
	)

	respData, sendResponse, err := Send(n.Client, n.BaseUrl, &sendRequest, n.SMSKey)

	return &smsclient.SendResult{
		ClientResponse: respData,
		Success:        sendResponse.StatusCode == 200,
	}, nil
}

var _ smsclient.RawClient = &SendCloudClient{}
