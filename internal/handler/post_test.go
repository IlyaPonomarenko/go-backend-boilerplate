package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api-consume-example/internal/repository"
	"api-consume-example/internal/service"
)

func testServer() http.Handler {
	postRepository := repository.NewMemoryPostRepository()
	postService := service.NewPostService(postRepository)
	postHandler := NewPostHandler(postService)
	mux := http.NewServeMux()
	postHandler.RegisterRoutes(mux)
	return mux
}

func TestCreateAndGetPost(t *testing.T) {
	server := testServer()
	request := httptest.NewRequest(http.MethodPost, "/api/posts", strings.NewReader(`{"userId":1,"title":"Hello","body":"World"}`))
	request.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()
	server.ServeHTTP(createResponse, request)

	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", createResponse.Code, http.StatusCreated)
	}

	var created map[string]any
	if err := json.NewDecoder(createResponse.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	if created["id"] != float64(1) {
		t.Fatalf("created ID = %v, want 1", created["id"])
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/api/posts/1", nil)
	getResponse := httptest.NewRecorder()
	server.ServeHTTP(getResponse, getRequest)
	if getResponse.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", getResponse.Code, http.StatusOK)
	}
}

func TestCreateRejectsMissingTitle(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/posts", strings.NewReader(`{"userId":1,"title":"  "}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	testServer().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestCreateRejectsNonJSON(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/posts", strings.NewReader("title=Hello"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	testServer().ServeHTTP(response, request)

	if response.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnsupportedMediaType)
	}
}
