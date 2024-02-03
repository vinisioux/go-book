package database

import (
	"database/sql"
	"go-book-api/src/config"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.ConnectionStringDatabase)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
