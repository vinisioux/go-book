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

func (repository Posts) Get(userId uint64) ([]models.Post, error) {
	rows, err := repository.db.Query(`
		select distinct p.*, u.nickname
		from posts p
		inner join users u
		on p.author_id = u.id
		inner join followers f
		on p.author_id = f.user_id
		where u.id = $1 or f.follower_id = $1
		order by p.created_at desc
	`, userId)
	if err != nil {
		return []models.Post{}, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		if err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorId,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNickname,
		); err != nil {
			return []models.Post{}, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (repository Posts) Update(postId uint64, post models.Post) error {
	statement, err := repository.db.Prepare("update posts set title = $1, content = $2 where id = $3")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(post.Title, post.Content, postId); err != nil {
		return err
	}

	return nil
}

func (repository Posts) Delete(postId uint64) error {
	statement, err := repository.db.Prepare("delete from posts where id = $1")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(postId); err != nil {
		return err
	}

	return nil
}

func (repository Posts) GetByUser(userId uint64) ([]models.Post, error) {
	rows, err := repository.db.Query(`
		select p.*, u.nickname
		from posts p inner join users u
		on p.author_id = u.id
		where p.author_id = $1
		order by p.created_at desc
	`, userId)
	if err != nil {
		return []models.Post{}, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		if err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorId,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNickname,
		); err != nil {
			return []models.Post{}, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (repository Posts) Like(postId uint64) error {
	statement, err := repository.db.Prepare("update posts set likes = likes + 1 where id = $1")
	if err != nil {
		return nil
	}
	defer statement.Close()

	if _, err := statement.Exec(postId); err != nil {
		return err
	}

	return nil
}

func (repository Posts) Unlike(postId uint64) error {
	statement, err := repository.db.Prepare(`
		update posts set likes =
		CASE 
			WHEN likes > 0 THEN likes - 1
			ELSE 0
		END
		where id = $1
	`)
	if err != nil {
		return nil
	}
	defer statement.Close()

	if _, err := statement.Exec(postId); err != nil {
		return err
	}

	return nil
}
