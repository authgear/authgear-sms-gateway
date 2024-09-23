package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/type_util"
)

type AccessYouClient struct {
	BaseUrl   string
	Client    *http.Client
	AccountNo string
	User      string
	Pwd       string
	Sender    string
	Logger    *slog.Logger
}

func NewAccessYouClient(
	baseUrl string,
	accountNo string,
	user string,
	pwd string,
	sender string,
	logger *slog.Logger,
) *AccessYouClient {
	if baseUrl == "" {
		baseUrl = "http://sms.accessyou-anyip.com"
	}
	return &AccessYouClient{
		BaseUrl:   baseUrl,
		Client:    &http.Client{},
		AccountNo: accountNo,
		User:      user,
		Pwd:       pwd,
		Sender:    sender,
		Logger:    logger,
	}
}

func (n *AccessYouClient) Send(options *SendOptions) (*SendResult, error) {
	// Access you phone number should have no +
	m1 := regexp.MustCompile(`[\+\-]+`)
	to := m1.ReplaceAllString(options.To, "")
	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"%v/sendsms.php?accountno=%v&pwd=%v&tid=1&phone=%v&a=%v&user=%v&from=%v",
			n.BaseUrl,
			n.AccountNo,
			n.Pwd,
			to,
			url.QueryEscape(options.Body),
			n.User,
			n.Sender),
		nil)
	req.Header.Set("Cookie", "dynamic=sms")

	resp, err := n.Client.Do(req)
	if err != nil {
		n.Logger.Error(fmt.Sprintf("%v", err))
		return nil, err
	}
	defer resp.Body.Close()

	n.Logger.Info("Attempt to parse")
	respData, err := io.ReadAll(resp.Body)

	// The response data is in format
	// "\ufeff{\"msg_status\":\"100\",\"msg_status_desc\":\"Successfully submitted message. \\u6267\\u884c\\u6210\\u529f\",\"phoneno\":\"852********\",\"msg_id\":852309279}"

	// Remove BOM token from resp json
	respData = bytes.Replace(respData, []byte("\ufeff"), []byte(""), -1)

	accessYouResponse := &AccessYouResponse{}
	err = json.Unmarshal(respData, &accessYouResponse)
	if err != nil {
		n.Logger.Error(fmt.Sprintf("Unmarshal error: %v", err))
		return nil, err
	}
	n.Logger.Info(fmt.Sprintf("%v", accessYouResponse))
	return &SendResult{
		ClientResponse: respData,
	}, nil
}

type AccessYouResponse struct {
	Status      string                         `json:"msg_status"`
	Description string                         `json:"msg_status_desc"`
	PhoneNo     type_util.SensitivePhoneNumber `json:"phoneno"`
}

var _ RawClient = &AccessYouClient{}
