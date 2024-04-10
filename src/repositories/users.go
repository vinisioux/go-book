package repositories

import (
	"database/sql"
	"fmt"
	"go-book-api/src/models"
)

type users struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *users {
	return &users{db}
}

func (repository users) Create(user models.User) (uint, error) {
	var userId uint
	err := repository.db.QueryRow(
		"insert into users (name, nickname, email, password) values ($1, $2, $3, $4) returning id",
		user.Name, user.Nickname, user.Email, user.Password,
	).Scan(&userId)

	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (repository users) Find(content string) ([]models.User, error) {
	content = fmt.Sprintf("%%%s%%", content) // %content%

	rows, err := repository.db.Query(
		"select id, name, nickname, email, created_at from users where name LIKE $1 or nickname LIKE $1",
		content,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (repository users) FindById(id uint64) (models.User, error) {
	rows, err := repository.db.Query(
		"select id, name, nickname, email, created_at from users where id = $1",
		id,
	)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	var user models.User

	if rows.Next() {
		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (repository users) Update(id uint64, user models.User) error {
	statement, err := repository.db.Prepare(
		"update users set name = $1, nickname = $2, email = $3 where id = $4",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(user.Name, user.Nickname, user.Email, id); err != nil {
		return err
	}

	return nil
}

func (repository users) Delete(id uint64) error {
	statement, err := repository.db.Prepare(
		"delete from users where id = $1",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(id); err != nil {
		return err
	}

	return nil
}

func (repository users) FindByEmail(email string) (models.User, error) {
	row, err := repository.db.Query("select id, password from users where email = $1", email)
	if err != nil {
		return models.User{}, err
	}
	defer row.Close()

	var user models.User

	if row.Next() {
		if err = row.Scan(&user.ID, &user.Password); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (repository users) Follow(userId, followerId uint64) error {
	statement, err := repository.db.Prepare(
		"insert into followers (user_id, follower_id) values ($1, $2) on conflict do nothing",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(userId, followerId); err != nil {
		return err
	}

	return nil
}

func (repository users) Unfollow(userId, followerId uint64) error {
	statement, err := repository.db.Prepare(
		"delete from followers where user_id = $1 and follower_id = $2",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(userId, followerId); err != nil {
		return err
	}

	return nil
}

func (repository users) GetFollowers(userId uint64) ([]models.User, error) {
	rows, err := repository.db.Query(`
		select
			u.id, u.name, u.nickname, u.email, u.created_at
		from users u inner join followers f 
		on u.id = f.follower_id
		where f.user_id = $1
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (repository users) GetFollowing(userId uint64) ([]models.User, error) {
	rows, err := repository.db.Query(`
		select
			u.id, u.name, u.nickname, u.email, u.created_at
		from users u inner join followers f 
		on u.id = f.user_id
		where f.follower_id = $1
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
