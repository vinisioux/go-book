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

func (repository Posts) GetById(postId uint64) (models.Post, error) {
	row, err := repository.db.Query(`
		select p.*, u.nickname
		from posts p inner join users u
		on p.author_id = u.id
		where p.id = $1
	`, postId)
	if err != nil {
		return models.Post{}, err
	}
	defer row.Close()

	var post models.Post

	if row.Next() {
		if err = row.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorId,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNickname,
		); err != nil {
			return models.Post{}, err
		}
	}

	return post, nil
}
