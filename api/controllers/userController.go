package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dmdinh22/go-blog/api/auth"
	"github.com/dmdinh22/go-blog/api/models"
	"github.com/dmdinh22/go-blog/api/responses"
	"github.com/dmdinh22/go-blog/api/utils/formaterror"
	"github.com/gorilla/mux"
)

// CreateUser godoc
// @Summary Creates a new user
// @Description Registers a new user to the DB
// @Tags users
// @Param login body models.User false "JSON request body for user"
// @Param ID query int false "user's id number"
// @Param Email query string true "user's email"
// @Param Username query string true "user's username"
// @Param Password query string true "user's password"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Router /api/users [post]
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("")

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	createdUser, err := user.CreateUser(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, createdUser.ID))
	responses.JSON(w, http.StatusCreated, createdUser)
}

// GetUsers godoc
// @Summary Get all users
// @Description Get details of all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Router /api/users [get]
func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	users, err := user.GetAllUsers(server.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

// GetUser godoc
// @Summary Get User By ID
// @Description Get details of a user by ID
// @Tags users
// @Param id path int true "User ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Router /api/users/{id} [get]
func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := models.User{}
	userRetrieved, err := user.GetUserById(server.DB, uint32(uid))

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, userRetrieved)
}

// Update User godoc
// @Summary Update User By ID
// @Description Update details of a user by ID
// @Tags users
// @Param id path int true "User ID"
// @Param user body models.User true "Update Request Body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Router /api/users/{id} [put]
func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenID, err := auth.ExtractTokenId(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	user.Prepare()
	err = user.Validate("update")

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedUser, err := user.UpdateUser(server.DB, uint32(uid))

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, updatedUser)
}

// Delete User godoc
// @Summary Delete User By ID
// @Description Delete details of a user by ID
// @Tags users
// @Param id path int true "User ID"
// @Accept  json
// @Produce  json
// @Success 200
// @Router /api/users/{id} [delete]
func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := models.User{}
	uid, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtractTokenId(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != 0 && tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = user.DeleteUser(server.DB, uint32(uid))

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Also delete the posts, likes and the comments that this user created if any:
	comment := models.Comment{}
	post := models.Post{}

	_, err = post.DeleteUserPosts(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	_, err = comment.DeleteUserComments(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"token": "User has been deleted",
	}

	responses.JSON(w, http.StatusOK, response)
}
