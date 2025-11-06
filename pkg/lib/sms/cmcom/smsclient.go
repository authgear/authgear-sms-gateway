package cmcom

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type CMCOMClient struct {
	Client *http.Client

	From         string
	ProductToken string

	Logger *slog.Logger
}

func (c *CMCOMClient) Send(ctx context.Context, options *smsclient.SendOptions) (*smsclient.SendResultSuccess, error) {
	return SendMessage(ctx, c.Client, c.Logger, c.ProductToken, c.From, string(options.To), options.Body, "")
}

func (c *CMCOMClient) ProviderType() string {
	return "cmcom"
}

var _ smsclient.RawClient = &CMCOMClient{}
