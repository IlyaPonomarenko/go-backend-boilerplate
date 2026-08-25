package repository

import (
	"context"
	"sync"

	"api-consume-example/internal/model"
)

// MemoryPostRepository stores posts in a map protected by a read/write lock.
// It is useful for learning and local development; its data disappears when
// the process stops.
type MemoryPostRepository struct {
	mu     sync.RWMutex
	nextID int
	posts  map[int]model.Post
}

// NewMemoryPostRepository creates an empty repository whose first ID is 1.
func NewMemoryPostRepository() *MemoryPostRepository {
	return &MemoryPostRepository{
		nextID: 1,
		posts:  make(map[int]model.Post),
	}
}

// List returns a snapshot, allowing callers to use the results after the read
// lock has been released.
func (repository *MemoryPostRepository) List(context.Context) ([]model.Post, error) {
	repository.mu.RLock()
	posts := make([]model.Post, 0, len(repository.posts))
	for _, post := range repository.posts {
		posts = append(posts, post)
	}
	repository.mu.RUnlock()
	return posts, nil
}

// Create assigns an ID and stores the new post atomically.
func (repository *MemoryPostRepository) Create(_ context.Context, input model.PostInput) (model.Post, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	post := model.Post{ID: repository.nextID, UserID: input.UserID, Title: input.Title, Body: input.Body}
	repository.posts[post.ID] = post
	repository.nextID++
	return post, nil
}

// Get returns a post and whether the requested ID exists.
func (repository *MemoryPostRepository) Get(_ context.Context, id int) (model.Post, bool, error) {
	repository.mu.RLock()
	post, found := repository.posts[id]
	repository.mu.RUnlock()
	return post, found, nil
}

// Update replaces an existing post while preserving its URL-selected ID.
func (repository *MemoryPostRepository) Update(_ context.Context, id int, input model.PostInput) (model.Post, bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	if _, found := repository.posts[id]; !found {
		return model.Post{}, false, nil
	}
	post := model.Post{ID: id, UserID: input.UserID, Title: input.Title, Body: input.Body}
	repository.posts[id] = post
	return post, true, nil
}

// Delete removes a post and reports whether it existed.
func (repository *MemoryPostRepository) Delete(_ context.Context, id int) (bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	if _, found := repository.posts[id]; !found {
		return false, nil
	}
	delete(repository.posts, id)
	return true, nil
}
