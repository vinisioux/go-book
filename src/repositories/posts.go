package repositories

import (
	"database/sql"
	"go-book-api/src/models"
)

type Posts struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *Posts {
	return &Posts{db}
}

func (repository Posts) Create(post models.Post) (uint64, error) {
	var postId uint64
	err := repository.db.QueryRow(
		"insert into posts (title, content, author_id) values ($1, $2, $3) returning id",
		post.Title, post.Content, post.AuthorId,
	).Scan(&postId)
	if err != nil {
		return 0, err
	}

	return postId, nil
}
