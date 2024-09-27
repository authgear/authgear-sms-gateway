package accessyou

import (
	"encoding/json"
)

type SendSMSResponse struct {
	MessageID   int    `json:"msg_id"`
	Status      string `json:"msg_status"`
	Description string `json:"msg_status_desc"`
	PhoneNo     string `json:"phoneno"`
}

func ParseSendSMSResponse(jsonData []byte) (*SendSMSResponse, error) {
	response := &SendSMSResponse{}
	err := json.Unmarshal(jsonData, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
