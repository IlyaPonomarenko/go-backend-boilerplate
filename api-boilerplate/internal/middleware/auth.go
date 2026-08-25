package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

// BearerAuth protects API routes with a constant-time comparison of the
// configured token. Health checks can remain outside this middleware.
func BearerAuth(expected string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			provided := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
			if provided == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) != 1 {
				writer.Header().Set("WWW-Authenticate", `Bearer realm="api"`)
				writeJSONError(writer, http.StatusUnauthorized, "unauthorized")
				return
			}
			next.ServeHTTP(writer, request)
		})
	}
}

func writeJSONError(writer http.ResponseWriter, status int, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_, _ = writer.Write([]byte(`{"error":"` + message + `"}` + "\n"))
}
