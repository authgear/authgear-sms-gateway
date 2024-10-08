package sendcloud

import (
	// nolint: gosec
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
	// nolint: gosec
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
	// According to the [doc](https://www.sendcloud.net/doc/sms/),
	// - The keys should be arranged alphabetically
	// - The values no need to be url encoded
	res := strings.Join([]string{
		fmt.Sprintf("msgType=%v", r.msgType),
		fmt.Sprintf("phone=%v", strings.Join(r.phone, ",")),
		fmt.Sprintf("sendRequestId=%v", r.sendRequestId),
		fmt.Sprintf("smsUser=%v", r.smsUser),
		fmt.Sprintf("templateId=%v", r.templateId),
	}, "&")

	if len(r.vars) != 0 {
		vars, _ := json.Marshal(r.vars)
		res = strings.Join([]string{
			res,
			fmt.Sprintf("vars=%v", string(vars)),
		}, "&")
	}
	return res
}

func (r *SendRequest) ToValues() url.Values {
	values := url.Values{}
	values.Set("msgType", r.msgType)
	values.Set("phone", strings.Join(r.phone, ","))
	values.Set("sendRequestId", r.sendRequestId)
	values.Set("smsUser", r.smsUser)
	values.Set("templateId", r.templateId)
	if len(r.vars) != 0 {
		vars, _ := json.Marshal(r.vars)
		values.Set("vars", string(vars))
	}
	return values
}

func (r *SendRequest) Sign(key string) string {
	signStr := fmt.Sprintf("%v&%v&%v", key, r.Presign(), key)
	// nolint: gosec
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
