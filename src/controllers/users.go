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
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
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

	if err = user.Prepare("register"); err != nil {
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
	user.ID, err = repository.Create(user)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	content := strings.ToLower(r.URL.Query().Get("content"))

	db, err := database.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)

	users, err := repository.Find(content)

	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId, err := strconv.ParseUint(params["userId"], 10, 64)

	if err != nil {
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

	user, err := repository.FindById(userId)

	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["userId"], 10, 64)

	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	tokenUserId, err := auth.GetUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	if userId != tokenUserId {
		responses.Err(w, http.StatusForbidden, errors.New("invalid userId"))
		return
	}

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

	if err = user.Prepare("update"); err != nil {
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
	if err = repository.Update(userId, user); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["userId"], 10, 64)

	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	tokenUserId, err := auth.GetUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	if userId != tokenUserId {
		responses.Err(w, http.StatusForbidden, errors.New("invalid userId"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)

	if err = repository.Delete(userId); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	followerId, err := auth.GetUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)

	userId, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if followerId == userId {
		responses.Err(w, http.StatusForbidden, errors.New("cannot follow yourself"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)

	if err = repository.Follow(userId, followerId); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	followerId, err := auth.GetUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)

	userId, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if followerId == userId {
		responses.Err(w, http.StatusForbidden, errors.New("cannot unfollow yourself"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	if err = repository.Unfollow(userId, followerId); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
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

	followers, err := repository.GetFollowers(userId)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, followers)
}

func GetFollowing(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
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

	users, err := repository.GetFollowing(userId)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	tokenUserId, err := auth.GetUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	if tokenUserId != userId {
		responses.Err(w, http.StatusForbidden, errors.New("cannot update password from another user"))
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	var password models.Password

	if err = json.Unmarshal(bodyRequest, &password); err != nil {
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

	storedPassword, err := repository.GetPassword(userId)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.CheckHash(storedPassword, password.CurrentPassword); err != nil {
		responses.Err(w, http.StatusUnauthorized, errors.New("wrong password"))
		return
	}

	hashedPassword, err := security.Hash(password.NewPassword)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = repository.UpdatePassword(userId, string(hashedPassword)); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
