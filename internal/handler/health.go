package handler

import "net/http"

// Health is a lightweight readiness check for monitoring or load balancers.
// It does not access the post repository.
func Health(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}
