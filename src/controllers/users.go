package controllers

import (
	"encoding/json"
	"fmt"
	"go-book-api/src/database"
	"go-book-api/src/models"
	"go-book-api/src/repositories"
	"io"
	"log"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	var user models.User
	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	repository := repositories.NewUserRepository(db)
	userId, err := repository.Create(user)
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(fmt.Sprintf("User created: %d", userId)))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get users"))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get user"))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update user"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user"))
}
