package api

import "net/http"


type HealthCheckHandler struct {}

func NewHealthCheckHandler () *HealthCheckHandler {
	return &HealthCheckHandler{}
}
func (hch *HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	SendJSON(w, map[string]string{"status": "ok"})
}
