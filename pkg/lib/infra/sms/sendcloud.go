package sms

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrMissingSendCloudConfiguration = errors.New("accessyou: configuration is missing")

type SendCloudMsgType int

const (
	SendCloudMsgTypeSMS SendCloudMsgType = 0
)

type SendCloudRequest struct {
	msgType       SendCloudMsgType
	phone         []string
	sendRequestId string
	smsUser       string
	templateId    int
	vars          map[string]interface{}
}

func NewSendCloudRequest(
	msgType SendCloudMsgType,
	phone []string,
	smsUser string,
	templateId int,
	vars map[string]interface{},
) *SendCloudRequest {
	s := &SendCloudRequest{
		msgType:    msgType,
		phone:      phone,
		smsUser:    smsUser,
		templateId: templateId,
		vars:       vars,
	}
	presign := s.Presign()
	h := md5.Sum([]byte(presign))
	sendRequestId := hex.EncodeToString(h[:])
	return &SendCloudRequest{
		msgType:       msgType,
		phone:         phone,
		sendRequestId: sendRequestId,
		smsUser:       smsUser,
		templateId:    templateId,
		vars:          vars,
	}
}

func (r *SendCloudRequest) Presign() string {
	vars, _ := json.Marshal(r.vars)
	return strings.Join([]string{
		fmt.Sprintf("msgType=%v", r.msgType),
		fmt.Sprintf("phone=%v", strings.Join(r.phone, ",")),
		fmt.Sprintf("sendRequestId=%v", r.sendRequestId),
		fmt.Sprintf("smsUser=%v", r.smsUser),
		fmt.Sprintf("templateId=%v", r.templateId),
		fmt.Sprintf("vars=%v", vars),
	}, "&")
}

func (r *SendCloudRequest) ToMap() map[string]interface{} {
	vars, _ := json.Marshal(r.vars)
	return map[string]interface{}{
		"msgType":       r.msgType,
		"phone":         strings.Join(r.phone, ","),
		"sendRequestId": r.sendRequestId,
		"smsUser":       r.smsUser,
		"templateId":    r.templateId,
		"vars":          vars,
	}
}

func Sign(key string, content string) string {
	signStr := fmt.Sprintf("%v&%v&%v", key, content, key)
	h := md5.Sum([]byte(signStr))
	return hex.EncodeToString(h[:])
}

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

func (n *SendCloudClient) Send(to string, body string) error {

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
