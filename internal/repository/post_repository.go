package repository

import (
	"api-consume-example/internal/model"
	"context"
)

// PostRepository describes the storage operations needed by the post service.
// A database implementation can satisfy this interface without changing the
// service or HTTP handler packages.
type PostRepository interface {
	List(context.Context) ([]model.Post, error)
	Create(context.Context, model.PostInput) (model.Post, error)
	Get(context.Context, int) (model.Post, bool, error)
	Update(context.Context, int, model.PostInput) (model.Post, bool, error)
	Delete(context.Context, int) (bool, error)
}
