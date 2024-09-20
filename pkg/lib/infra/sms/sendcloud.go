package sms

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	. "github.com/authgear/authgear-sms-gateway/pkg/lib/infra/sms/sendcloud"
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
	Name             string
	BaseUrl          string
	Client           *http.Client
	SMSUser          string
	SMSKey           string
	TemplateResolver ISendCloudTemplateResolver
	Logger           *slog.Logger
}

func NewSendCloudClient(
	name string,
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
		Name:             name,
		BaseUrl:          baseUrl,
		Client:           &http.Client{},
		SMSUser:          smsUser,
		SMSKey:           smsKey,
		TemplateResolver: templateResolver,
		Logger:           logger,
	}
}

func (n *SendCloudClient) GetName() string {
	return n.Name
}

func (n *SendCloudClient) Send(
	to string,
	body string,
	templateName string,
	languageTag string,
	templateVariables *TemplateVariables,
) (ClientResponse, error) {
	template, err := n.TemplateResolver.Resolve(templateName, languageTag)
	if err != nil {
		return ClientResponse{}, err
	}
	sendCloudRequest := NewSendCloudRequest(
		string(template.TemplateMsgType),
		[]string{
			to,
		},
		n.SMSUser,
		string(template.TemplateID),
		makeVarsFromTemplateVariables(templateVariables),
	)

	n.Logger.Debug(fmt.Sprintf("Presign: %v", sendCloudRequest.Presign()))
	values := sendCloudRequest.ToValues()
	values.Set("signature", sendCloudRequest.Sign(n.SMSKey))

	data := values.Encode()
	n.Logger.Debug(fmt.Sprintf("data: %v", data))

	req, _ := http.NewRequest("POST", fmt.Sprintf("%v/smsapi/send", n.BaseUrl), strings.NewReader(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := n.Client.Do(req)

	if err != nil {
		n.Logger.Error(fmt.Sprintf("Client.Do error: %v", err))
		return ClientResponse{}, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	n.Logger.Error(fmt.Sprintf("resp: %v", string(respData)))

	return ClientResponse(respData), nil
}

var _ RawClient = &SendCloudClient{}
