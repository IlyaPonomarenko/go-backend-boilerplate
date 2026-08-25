package service

import (
	"context"
	"errors"
	"strings"

	"api-consume-example/internal/model"
	"api-consume-example/internal/repository"
)

// ErrInvalidTitle means a post does not contain the required title.
var ErrInvalidTitle = errors.New("title is required")

// PostService contains post business rules and coordinates persistence. It
// does not know whether the repository uses memory, SQL, or another backend.
type PostService struct {
	repository repository.PostRepository
}

// NewPostService creates a post service with its storage dependency.
func NewPostService(postRepository repository.PostRepository) *PostService {
	return &PostService{repository: postRepository}
}

// List returns all posts currently available from the repository.
func (service *PostService) List(ctx context.Context) ([]model.Post, error) {
	return service.repository.List(ctx)
}

// Create validates input before asking the repository to assign an ID.
func (service *PostService) Create(ctx context.Context, input model.PostInput) (model.Post, error) {
	if strings.TrimSpace(input.Title) == "" {
		return model.Post{}, ErrInvalidTitle
	}
	return service.repository.Create(ctx, input)
}

// Get retrieves one post. The boolean is false when the ID does not exist.
func (service *PostService) Get(ctx context.Context, id int) (model.Post, bool, error) {
	return service.repository.Get(ctx, id)
}

// Update validates input and replaces an existing post.
func (service *PostService) Update(ctx context.Context, id int, input model.PostInput) (model.Post, bool, error) {
	if strings.TrimSpace(input.Title) == "" {
		return model.Post{}, false, ErrInvalidTitle
	}
	return service.repository.Update(ctx, id, input)
}

// Delete removes one post and reports whether it existed.
func (service *PostService) Delete(ctx context.Context, id int) (bool, error) {
	return service.repository.Delete(ctx, id)
}
