package sendcloud

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type SendCloudRequest struct {
	msgType       string
	phone         []string
	sendRequestId string
	smsUser       string
	templateId    string
	vars          map[string]interface{}
}

func NewSendCloudRequest(
	msgType string,
	phone []string,
	smsUser string,
	templateId string,
	vars map[string]interface{},
) SendCloudRequest {
	s := SendCloudRequest{
		msgType:    msgType,
		phone:      phone,
		smsUser:    smsUser,
		templateId: templateId,
		vars:       vars,
	}
	presign := s.Presign()
	h := md5.Sum([]byte(presign))
	sendRequestId := hex.EncodeToString(h[:])
	return SendCloudRequest{
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
		fmt.Sprintf("vars=%v", string(vars)),
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

func (r *SendCloudRequest) ToValues() url.Values {
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

func (r *SendCloudRequest) Sign(key string) string {
	signStr := fmt.Sprintf("%v&%v&%v", key, r.Presign(), key)
	h := md5.Sum([]byte(signStr))
	return hex.EncodeToString(h[:])
}
