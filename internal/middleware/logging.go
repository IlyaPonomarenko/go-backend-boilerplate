package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Logging measures each request and logs it after the wrapped handler
// completes. Middleware is a good home for behavior shared by every route.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		started := time.Now()
		next.ServeHTTP(writer, request)
		slog.Info("request complete", "method", request.Method, "path", request.URL.Path, "duration_ms", time.Since(started).Milliseconds())
	})
}
