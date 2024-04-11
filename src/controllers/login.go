package controllers

import (
	"encoding/json"
	"errors"
	"go-book-api/src/auth"
	"go-book-api/src/database"
	"go-book-api/src/models"
	"go-book-api/src/repositories"
	"go-book-api/src/responses"
	"go-book-api/src/security"
	"io"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	userFound, err := repository.FindByEmail(user.Email)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.CheckHash(userFound.Password, user.Password); err != nil {
		responses.Err(w, http.StatusUnauthorized, errors.New("wrong credentials"))
		return
	}

	token, err := auth.CreateToken(userFound.ID)

	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	w.Write([]byte(token))
}
