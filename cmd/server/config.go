package main

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ListenAddr                   string `envconfig:"LISTEN_ADDR" default:"0.0.0.0:8091"`
	SMSServiceProviderConfigPath string `envconfig:"SMS_SERVICE_PROVIDER_CONFIG_PATH"`
}

func LoadConfigFromEnv() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
