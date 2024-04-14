package controllers

import (
	"encoding/json"
	"go-book-api/src/auth"
	"go-book-api/src/database"
	"go-book-api/src/models"
	"go-book-api/src/repositories"
	"go-book-api/src/responses"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post models.Post

	if err = json.Unmarshal(bodyRequest, &post); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	post.AuthorId = userId

	if err = post.Prepare(); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewPostRepository(db)

	post.ID, err = repository.Create(post)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, post)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {

}

func GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	postId, err := strconv.ParseUint(params["postId"], 10, 64)
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

	repository := repositories.NewPostRepository(db)

	post, err := repository.GetById(postId)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, post)

}

func UpdatePost(w http.ResponseWriter, r *http.Request) {

}

func DeletePost(w http.ResponseWriter, r *http.Request) {

}
