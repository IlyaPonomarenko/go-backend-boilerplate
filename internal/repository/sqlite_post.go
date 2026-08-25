package repository

import (
	"context"
	"database/sql"

	"api-consume-example/internal/model"
)

// SQLitePostRepository persists posts in SQLite. The schema is initialized by
// the constructor so a fresh deployment can start without manual SQL steps.
type SQLitePostRepository struct {
	db *sql.DB
}

// NewSQLitePostRepository opens the database and creates the posts table.
func NewSQLitePostRepository(ctx context.Context, db *sql.DB) (*SQLitePostRepository, error) {
	repository := &SQLitePostRepository{db: db}
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			body TEXT NOT NULL
		)`)
	if err != nil {
		return nil, err
	}
	return repository, nil
}

func (repository *SQLitePostRepository) List(ctx context.Context) ([]model.Post, error) {
	rows, err := repository.db.QueryContext(ctx, `SELECT id, user_id, title, body FROM posts ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]model.Post, 0)
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Body); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

func (repository *SQLitePostRepository) Create(ctx context.Context, input model.PostInput) (model.Post, error) {
	result, err := repository.db.ExecContext(ctx, `INSERT INTO posts (user_id, title, body) VALUES (?, ?, ?)`, input.UserID, input.Title, input.Body)
	if err != nil {
		return model.Post{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.Post{}, err
	}
	return model.Post{ID: int(id), UserID: input.UserID, Title: input.Title, Body: input.Body}, nil
}

func (repository *SQLitePostRepository) Get(ctx context.Context, id int) (model.Post, bool, error) {
	var post model.Post
	err := repository.db.QueryRowContext(ctx, `SELECT id, user_id, title, body FROM posts WHERE id = ?`, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Body)
	if err == sql.ErrNoRows {
		return model.Post{}, false, nil
	}
	if err != nil {
		return model.Post{}, false, err
	}
	return post, true, nil
}

func (repository *SQLitePostRepository) Update(ctx context.Context, id int, input model.PostInput) (model.Post, bool, error) {
	result, err := repository.db.ExecContext(ctx, `UPDATE posts SET user_id = ?, title = ?, body = ? WHERE id = ?`, input.UserID, input.Title, input.Body, id)
	if err != nil {
		return model.Post{}, false, err
	}
	count, err := result.RowsAffected()
	if err != nil || count == 0 {
		return model.Post{}, count > 0, err
	}
	return model.Post{ID: id, UserID: input.UserID, Title: input.Title, Body: input.Body}, true, nil
}

func (repository *SQLitePostRepository) Delete(ctx context.Context, id int) (bool, error) {
	result, err := repository.db.ExecContext(ctx, `DELETE FROM posts WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count > 0, err
}
