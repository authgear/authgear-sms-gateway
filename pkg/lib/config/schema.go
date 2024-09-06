package config

import (
	"github.com/authgear/authgear-server/pkg/util/validation"
)

var SMSProviderConfigSchema = validation.NewMultipartSchema("SMSProviderConfig")

func init() {
	SMSProviderConfigSchema.Instantiate()
}
