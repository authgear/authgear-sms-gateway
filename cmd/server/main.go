package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"github.com/authgear/authgear-sms-gateway/pkg/handler"
)

func main() {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("failed to load .env file: %s", err)
	}
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
