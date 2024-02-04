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
