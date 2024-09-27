package accessyou

import (
	"encoding/json"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/type_util"
)

type SendSMSResponse struct {
	MessageID   int                            `json:"msg_id"`
	Status      string                         `json:"msg_status"`
	Description string                         `json:"msg_status_desc"`
	PhoneNo     type_util.SensitivePhoneNumber `json:"phoneno"`
}

func ParseSendSMSResponse(jsonData []byte) (*SendSMSResponse, error) {
	response := &SendSMSResponse{}
	err := json.Unmarshal(jsonData, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
