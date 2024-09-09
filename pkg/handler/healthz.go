package handler

import (
	"fmt"
	"net/http"
)

type HealthzHandler struct {
}

func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

func (p *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
