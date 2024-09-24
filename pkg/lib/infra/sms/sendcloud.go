package sms

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	. "github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms/sendcloud"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms/sendcloud/apis"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms/sendcloud/models"
)

var ErrMissingSendCloudConfiguration = errors.New("accessyou: configuration is missing")

func makeVarsFromTemplateVariables(variables *TemplateVariables) map[string]interface{} {
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
		BaseUrl:          baseUrl,
		Client:           &http.Client{},
		SMSUser:          smsUser,
		SMSKey:           smsKey,
		TemplateResolver: templateResolver,
		Logger:           logger,
	}
}

func (n *SendCloudClient) Send(options *SendOptions) (*SendResult, error) {
	template, err := n.TemplateResolver.Resolve(options.TemplateName, options.LanguageTag)
	if err != nil {
		return nil, err
	}
	sendRequest := models.NewSendRequest(
		string(template.TemplateMsgType),
		[]string{
			options.To,
		},
		n.SMSUser,
		string(template.TemplateID),
		makeVarsFromTemplateVariables(options.TemplateVariables),
	)

	respData, sendResponse, err := apis.Send(n.Client, n.BaseUrl, &sendRequest, n.SMSKey)

	return &SendResult{
		ClientResponse: respData,
		Success:        sendResponse.StatusCode == 200,
	}, nil
}

var _ RawClient = &SendCloudClient{}
