package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBearerAuth(t *testing.T) {
	protected := BearerAuth("secret")(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	}))

	tests := []struct {
		name          string
		authorization string
		status        int
	}{
		{name: "missing token", status: http.StatusUnauthorized},
		{name: "wrong token", authorization: "Bearer wrong", status: http.StatusUnauthorized},
		{name: "valid token", authorization: "Bearer secret", status: http.StatusNoContent},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
			request.Header.Set("Authorization", test.authorization)
			response := httptest.NewRecorder()
			protected.ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d", response.Code, test.status)
			}
		})
	}
}
