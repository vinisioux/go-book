package models

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID        uint      `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nickname  string    `json:"nickname,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (user *User) Prepare(step string) error {
	if err := user.validar(step); err != nil {
		return err
	}

	user.formatFields()
	return nil
}

func (user *User) validar(step string) error {
	if user.Name == "" {
		return errors.New("name is required")
	}
	if user.Nickname == "" {
		return errors.New("nickname is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if step == "register" && user.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (user *User) formatFields() {
	user.Name = strings.TrimSpace(user.Name)
	user.Nickname = strings.TrimSpace(user.Nickname)
	user.Email = strings.TrimSpace(user.Email)
}
