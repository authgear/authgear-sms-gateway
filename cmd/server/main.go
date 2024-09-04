package main

import (
	"net/http"
	"time"

	"github.com/authgear/authgear-sms-gateway/pkg/handler"
)

func main() {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		panic(err)
	}

	http.Handle("/healthz", &handler.HealthzHandler{})

	server := &http.Server{
		Addr:              cfg.ListenAddr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
