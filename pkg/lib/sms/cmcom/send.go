package cmcom

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/api"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sensitive"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

const CM_API_ENDPOINT = "https://gw.cmtelecom.com/v1.0/message"

type MessageBody struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type To struct {
	Number string `json:"number"`
}

type Msg struct {
	From      string      `json:"from"`
	To        []To        `json:"to"`
	Body      MessageBody `json:"body"`
	Reference string      `json:"reference,omitempty"`
}

type Messages struct {
	Msg []Msg `json:"msg"`
}

type SendRequest struct {
	Messages Messages `json:"messages"`
}

type MessageResponse struct {
	To               string  `json:"to"`
	Status           string  `json:"status"`
	Reference        string  `json:"reference,omitempty"`
	Parts            int     `json:"parts"`
	MessageDetails   *string `json:"messageDetails,omitempty"`
	MessageErrorCode int     `json:"messageErrorCode"`
}

type SendResponse struct {
	Details   string            `json:"details"`
	ErrorCode int               `json:"errorCode"`
	Messages  []MessageResponse `json:"messages"`
}

type SendMessageOptions struct {
	ProductToken string
	From         string
	To           string
	Content      string
}

// SendMessage sends a message via CM.com API and returns the response or an error
func SendMessage(
	ctx context.Context,
	httpClient *http.Client,
	logger *slog.Logger,
	options *SendMessageOptions,
) (*smsclient.SendResultSuccess, error) {
	reqBody := SendRequest{
		Messages: Messages{
			Msg: []Msg{
				{
					From: options.From,
					To:   []To{{Number: options.To}},
					Body: MessageBody{
						Type:    "auto",
						Content: options.Content,
					},
				},
			},
		},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", CM_API_ENDPOINT, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CM-PRODUCTTOKEN", options.ProductToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		err = sensitive.RedactHTTPClientError(err)
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			err = errors.Join(err, &smsclient.SendResultError{
				DumpedResponse: nil,
				Code:           api.CodeTimeout,
			})
		}
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	dumpedResponse, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(
			err,
			&smsclient.SendResultError{
				DumpedResponse: dumpedResponse,
			},
		)
	}

	// Parse response JSON
	var sendResp SendResponse
	if err := json.Unmarshal(respData, &sendResp); err != nil {
		// Failed to parse response JSON
		sendErr := &smsclient.SendResultError{
			DumpedResponse: dumpedResponse,
		}
		var jsonSyntaxErr *json.SyntaxError
		if errors.As(err, &jsonSyntaxErr) {
			sendErr.Code = api.CodeUnknownError
		}
		return nil, errors.Join(
			err,
			sendErr,
		)
	}

	// Handle top‐level errorCode
	// Error code 201 means the request has errors in its messages
	// Need to check per-message errors
	if sendResp.ErrorCode != 0 && sendResp.ErrorCode != 201 {
		logger.ErrorContext(ctx, "cmcom error response",
			"errorCode", sendResp.ErrorCode,
			"details", sendResp.Details,
		)
		return nil, MakeError(sendResp.ErrorCode, dumpedResponse)
	}

	// Check per-message errors
	if len(sendResp.Messages) == 0 {
		// Empty messages in response is unexpected
		return nil, &smsclient.SendResultError{
			Code:           api.CodeUnknownError,
			DumpedResponse: dumpedResponse,
		}
	}

	msgResp := sendResp.Messages[0]
	logger.InfoContext(ctx, "cmcom message response",
		"status", msgResp.Status,
		"messageDetails", msgResp.MessageDetails,
		"messageErrorCode", msgResp.MessageErrorCode,
	)

	if msgResp.MessageErrorCode != 0 {
		return nil, MakeError(msgResp.MessageErrorCode, dumpedResponse)
	}

	return &smsclient.SendResultSuccess{
		DumpedResponse: dumpedResponse,
	}, nil
}

func MakeError(providerErrorCode int, dumpedResponse []byte) *smsclient.SendResultError {
	err := &smsclient.SendResultError{
		DumpedResponse:    dumpedResponse,
		ProviderErrorCode: fmt.Sprintf("%d", providerErrorCode),
	}

	// https://developers.cm.com/messaging/docs/responses-errors-json
	switch providerErrorCode {
	case 999: // Unknown error, please contact CM support
		err.Code = api.CodeUnknownError
	case 101: // Authentication of the request failed
		err.Code = api.CodeAuthenticationFailed
	case 102: // The account using this authentication has insufficient balance
		err.Code = api.CodeDeliveryRejected
	case 103: // The product token is incorrect
		err.Code = api.CodeAuthenticationFailed
	case 201: // This request has one or more errors in its messages. Some or all messages have not been sent. See MSGs for details
		err.Code = api.CodeDeliveryRejected
	case 202: // This request is malformed, please confirm the JSON and that the correct data types are used
		fallthrough
	case 203: // The request's MSG array is incorrect
		fallthrough
	case 301: // This MSG has an invalid From field (per msg)
		err.Code = api.CodeUnknownError
	case 302: // This MSG has an invalid To field (per msg)
		fallthrough
	case 303: // This MSG has an invalid Phone Number in the To field (per msg,)
		err.Code = api.CodeInvalidPhoneNumber
	case 304: // This MSG has an invalid Body field (per msg)
		fallthrough
	case 305: // This MSG has an invalid field. Please confirm with the documentation (per msg)
		fallthrough
	case 307: // This MSG exceeds the maximum size (per msg)
		fallthrough
	case 401: // Message has been spam filtered
		fallthrough
	case 402: // Message has been blacklisted
		fallthrough
	case 403: // Message has been rejected
		err.Code = api.CodeDeliveryRejected
	case 500: // An internal error has occurred
		err.Code = api.CodeUnknownError
	}

	return err
}
