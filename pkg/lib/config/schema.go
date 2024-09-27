package config

import (
	"github.com/authgear/authgear-server/pkg/util/validation"
)

var RootSchema = validation.NewMultipartSchema("RootConfig")

func init() {
	RootSchema.Instantiate()
}
