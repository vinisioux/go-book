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
		SELECT 
			pp.*, 
			COALESCE(l.likes, 0) AS likes
		FROM
			(SELECT p.*, u.nickname 
			FROM posts p
			INNER JOIN users u ON p.author_id = u.id
			WHERE p.id = $1) pp
		LEFT JOIN
			(SELECT post_id, COUNT(*) AS likes
			FROM likes
			WHERE post_id = $1
			GROUP BY post_id) l ON pp.id = l.post_id
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
			&post.CreatedAt,
			&post.AuthorNickname,
			&post.Likes,
		); err != nil {
			return models.Post{}, err
		}
	}

	return post, nil
}

func (repository Posts) Get(userId uint64) ([]models.Post, error) {
	rows, err := repository.db.Query(`
		SELECT 
			pp.*, 
			COALESCE(l.likes, 0) AS likes
		FROM
			(SELECT DISTINCT p.*, u.nickname 
			FROM posts p
			INNER JOIN users u ON p.author_id = u.id
			INNER JOIN followers f ON p.author_id = f.user_id
			WHERE u.id = $1 OR f.follower_id = $1) pp
		LEFT JOIN
			(SELECT post_id, COUNT(*) AS likes
			FROM likes
			WHERE user_id = $1
			GROUP BY post_id) l ON pp.id = l.post_id
		ORDER BY
			pp.created_at DESC
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
			&post.CreatedAt,
			&post.AuthorNickname,
			&post.Likes,
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
		SELECT 
			pp.*, 
			COALESCE(l.likes, 0) AS likes
		FROM
			(SELECT p.*, u.nickname 
			FROM posts p
			INNER JOIN users u ON p.author_id = u.id
			WHERE p.author_id = $1) pp
		LEFT JOIN
			(SELECT post_id, COUNT(*) AS likes
			FROM likes
			GROUP BY post_id) l ON pp.id = l.post_id
		ORDER BY 
			pp.created_at DESC;
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
			&post.CreatedAt,
			&post.AuthorNickname,
			&post.Likes,
		); err != nil {
			return []models.Post{}, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (repository Posts) Like(postId, userId uint64) error {
	statement, err := repository.db.Prepare(
		"insert into likes (post_id, user_id) values ($1, $2) on conflict (post_id, user_id) do nothing",
	)
	if err != nil {
		return nil
	}
	defer statement.Close()

	if _, err := statement.Exec(postId, userId); err != nil {
		return err
	}

	return nil
}

func (repository Posts) Unlike(postId, userId uint64) error {
	statement, err := repository.db.Prepare(
		"delete from likes where post_id = $1 and user_id = $2",
	)
	if err != nil {
		return nil
	}
	defer statement.Close()

	if _, err := statement.Exec(postId, userId); err != nil {
		return err
	}

	return nil
}
