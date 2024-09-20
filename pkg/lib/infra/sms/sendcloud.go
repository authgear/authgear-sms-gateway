package sms

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrMissingSendCloudConfiguration = errors.New("accessyou: configuration is missing")

type SendCloudClient struct {
	Name    string
	BaseUrl string
	Client  *http.Client
	SMSUser string
	SMSKey  string
}

func NewSendCloudClient(name string, baseUrl string, smsUser string, smsKey string) *SendCloudClient {
	if baseUrl == "" {
		baseUrl = "https://api.sendcloud.net"
	}
	return &SendCloudClient{
		Name:    name,
		BaseUrl: baseUrl,
		Client:  &http.Client{},
		SMSUser: smsUser,
		SMSKey:  smsKey,
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
) error {

	req, _ := http.NewRequest("POST", fmt.Sprintf("%v/smsapi/send", n.BaseUrl), nil)
	req.Header.Set("Cookie", "dynamic=sms")
	resp, err := n.Client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

var _ RawClient = &SendCloudClient{}
