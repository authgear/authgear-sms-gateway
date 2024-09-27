package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/authgear/authgear-sms-gateway/pkg/handler"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/config"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/logger"
	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms"
)

func main() {
	logger := logger.NewLogger()

	err := godotenv.Load()
	if errors.Is(err, os.ErrNotExist) {
		logger.Warn("skip loading .env as it is absent")
	} else if err != nil {
		panic(err)
	}

	// Set timeout to avoid indefinite waiting time.
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	envCfg, err := LoadEnvConfigFromEnv()
	if err != nil {
		panic(err)
	}

	configYAML, err := os.ReadFile(envCfg.ConfigPath)
	if err != nil {
		panic(err)
	}
	cfg, err := config.ParseRootConfigFromYAML(configYAML)
	if err != nil {
		panic(err)
	}
	smsClientMap := sms.NewSMSClientMap(cfg, httpClient, logger)
	smsService := &sms.SMSService{
		Logger:       logger,
		RootConfig:   cfg,
		SMSClientMap: smsClientMap,
	}

	http.Handle("/healthz", &handler.HealthzHandler{})
	http.Handle("/send", &handler.SendHandler{
		Logger:     logger,
		SMSService: smsService,
	})

	server := &http.Server{
		Addr:              envCfg.ListenAddr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	logger.Info(fmt.Sprintf("Server running at %v", envCfg.ListenAddr))
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
