package main

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ListenAddr string `envconfig:"LISTEN_ADDR" default:"0.0.0.0:8091"`
}

func LoadConfigFromEnv() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
