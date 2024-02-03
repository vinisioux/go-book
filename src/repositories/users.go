package repositories

import (
	"database/sql"
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
