package handler

import (
	"fmt"
	"net/http"
)

type HealthzHandler struct{}

func (p *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "OK")
}
