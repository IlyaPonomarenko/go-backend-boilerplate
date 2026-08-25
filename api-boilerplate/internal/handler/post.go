package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"api-consume-example/internal/model"
	"api-consume-example/internal/service"
)

// PostHandler translates HTTP requests into service calls. Business rules
// remain in PostService so they can be tested without an HTTP server.
type PostHandler struct {
	service *service.PostService
}

// NewPostHandler creates an HTTP handler with its service dependency.
func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{service: postService}
}

// RegisterRoutes attaches all post endpoints to the application's multiplexer.
func (handler *PostHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/posts", handler.list)
	mux.HandleFunc("POST /api/posts", handler.create)
	mux.HandleFunc("GET /api/posts/{id}", handler.get)
	mux.HandleFunc("PUT /api/posts/{id}", handler.update)
	mux.HandleFunc("DELETE /api/posts/{id}", handler.delete)
}

// list returns a JSON array with HTTP 200 OK.
func (handler *PostHandler) list(writer http.ResponseWriter, request *http.Request) {
	posts, err := handler.service.List(request.Context())
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "could not load posts")
		return
	}
	writeJSON(writer, http.StatusOK, posts)
}

// create decodes input, validates it through the service, and returns HTTP
// 201 Created with the new resource.
func (handler *PostHandler) create(writer http.ResponseWriter, request *http.Request) {
	var input model.PostInput
	if !decodeJSON(writer, request, &input) {
		return
	}

	post, err := handler.service.Create(request.Context(), input)
	if errors.Is(err, service.ErrInvalidTitle) {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "could not create post")
		return
	}
	writeJSON(writer, http.StatusCreated, post)
}

// get parses the route ID and returns HTTP 404 when the service cannot find it.
func (handler *PostHandler) get(writer http.ResponseWriter, request *http.Request) {
	id, ok := postID(writer, request)
	if !ok {
		return
	}

	post, found, err := handler.service.Get(request.Context(), id)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "could not load post")
		return
	}
	if !found {
		writeError(writer, http.StatusNotFound, "post not found")
		return
	}
	writeJSON(writer, http.StatusOK, post)
}

// update replaces an existing post using the ID from the URL.
func (handler *PostHandler) update(writer http.ResponseWriter, request *http.Request) {
	id, ok := postID(writer, request)
	if !ok {
		return
	}

	var input model.PostInput
	if !decodeJSON(writer, request, &input) {
		return
	}
	post, found, err := handler.service.Update(request.Context(), id, input)
	if errors.Is(err, service.ErrInvalidTitle) {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "could not update post")
		return
	}
	if !found {
		writeError(writer, http.StatusNotFound, "post not found")
		return
	}
	writeJSON(writer, http.StatusOK, post)
}

// delete removes an existing post and returns HTTP 204 No Content.
func (handler *PostHandler) delete(writer http.ResponseWriter, request *http.Request) {
	id, ok := postID(writer, request)
	if !ok {
		return
	}
	found, err := handler.service.Delete(request.Context(), id)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "could not delete post")
		return
	}
	if !found {
		writeError(writer, http.StatusNotFound, "post not found")
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

// postID centralizes validation for every endpoint with an {id} parameter.
func postID(writer http.ResponseWriter, request *http.Request) (int, bool) {
	id, err := strconv.Atoi(request.PathValue("id"))
	if err != nil || id < 1 {
		writeError(writer, http.StatusBadRequest, "id must be a positive integer")
		return 0, false
	}
	return id, true
}

// decodeJSON accepts only JSON request bodies and reports malformed input to
// the client. A false result means an error response has already been sent.
func decodeJSON(writer http.ResponseWriter, request *http.Request, target any) bool {
	if request.Header.Get("Content-Type") != "application/json" {
		writeError(writer, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return false
	}
	request.Body = http.MaxBytesReader(writer, request.Body, 1<<20)
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(target); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeError(writer, http.StatusBadRequest, "request body must contain one JSON object")
		return false
	}
	return true
}

// writeJSON keeps content type and status handling consistent for all JSON
// responses.
func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

// writeError gives every client-facing error the same JSON shape.
func writeError(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, map[string]string{"error": message})
}
