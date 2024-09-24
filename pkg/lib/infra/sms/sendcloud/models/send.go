package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type SendRequest struct {
	msgType       string
	phone         []string
	sendRequestId string
	smsUser       string
	templateId    string
	vars          map[string]interface{}
}

func NewSendRequest(
	msgType string,
	phone []string,
	smsUser string,
	templateId string,
	vars map[string]interface{},
) SendRequest {
	s := SendRequest{
		msgType:    msgType,
		phone:      phone,
		smsUser:    smsUser,
		templateId: templateId,
		vars:       vars,
	}
	presign := s.Presign()
	h := md5.Sum([]byte(presign))
	sendRequestId := hex.EncodeToString(h[:])
	return SendRequest{
		msgType:       msgType,
		phone:         phone,
		sendRequestId: sendRequestId,
		smsUser:       smsUser,
		templateId:    templateId,
		vars:          vars,
	}
}

func (r *SendRequest) Presign() string {
	vars, _ := json.Marshal(r.vars)

	// According to the [doc](https://www.sendcloud.net/doc/sms/),
	// - The keys should be arranged alphabetically
	// - The values no need to be url encoded
	return strings.Join([]string{
		fmt.Sprintf("msgType=%v", r.msgType),
		fmt.Sprintf("phone=%v", strings.Join(r.phone, ",")),
		fmt.Sprintf("sendRequestId=%v", r.sendRequestId),
		fmt.Sprintf("smsUser=%v", r.smsUser),
		fmt.Sprintf("templateId=%v", r.templateId),
		fmt.Sprintf("vars=%v", string(vars)),
	}, "&")
}

func (r *SendRequest) ToMap() map[string]interface{} {
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

func (r *SendRequest) ToValues() url.Values {
	vars, _ := json.Marshal(r.vars)
	values := url.Values{}
	values.Set("msgType", fmt.Sprintf("%v", r.msgType))
	values.Set("phone", strings.Join(r.phone, ","))
	values.Set("sendRequestId", r.sendRequestId)
	values.Set("smsUser", r.smsUser)
	values.Set("templateId", fmt.Sprintf("%v", r.templateId))
	values.Set("vars", string(vars))
	return values
}

func (r *SendRequest) Sign(key string) string {
	signStr := fmt.Sprintf("%v&%v&%v", key, r.Presign(), key)
	h := md5.Sum([]byte(signStr))
	return hex.EncodeToString(h[:])
}

type SendResponseInfo struct {
	SuccessCount int      `json:"successCount,omitempty"`
	SMSIDs       []string `json:"smsIds,omitempty"`
}

type SendResponse struct {
	Result     bool              `json:"result,omitempty"`
	StatusCode int               `json:"statusCode,omitempty"`
	Message    string            `json:"message,omitempty"`
	Info       *SendResponseInfo `json:"info"`
}

func ParseSendResponse(jsonData []byte) (*SendResponse, error) {
	response := &SendResponse{}
	err := json.Unmarshal(jsonData, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
