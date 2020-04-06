package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dmdinh22/go-blog/api/auth"
	"github.com/dmdinh22/go-blog/api/models"
	"github.com/dmdinh22/go-blog/api/responses"
	"github.com/dmdinh22/go-blog/api/utils/formaterror"
	"github.com/gorilla/mux"
)

// CreateComment godoc
// @Summary Creates a new comment
// @Description Creates a new comment for the user logged in
// @Tags comments
// @Param id query int false "comment's id "
// @Param userId query int true "comment's user Id"
// @Param postId query int true "comment's post Id"
// @Param body query string true "comment's content"
// @Param AuthorID query string true "id of user that created this comment"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Comment
// @Router /api/comments [post]
func (server *Server) CreateComment(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	comment := models.Comment{}
	err = json.Unmarshal(body, &comment)

	// validate token
	uid, err := auth.ExtractTokenId(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	// check user exists
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	// check post exists
	post := models.Post{}
	pid := comment.PostID
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)

		return
	}

	comment.Prepare()
	err = comment.Validate()

	// enter the userid and the postid. The comment body is automatically passed
	comment.UserID = uid
	comment.PostID = pid

	comment.Prepare()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	commentCreated, err := comment.AddComment(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	responses.JSON(w, http.StatusCreated, commentCreated)
}

// Delete Comment godoc
// @Summary Delete Comment By ID
// @Description Delete details of a Comment by ID
// @Tags comments
// @Param id path int true "Comment ID"
// @Accept  json
// @Produce  json
// @Success 200
// @Router /api/comments/{id} [delete]
func (server *Server) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Validate user
	uid, err := auth.ExtractTokenId(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	// Validate existence of comment
	comment := models.Comment{}
	err = server.DB.Debug().Model(models.Comment{}).Where("id = ?", cid).Take(&comment).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	// Validate that the comment belongs to this user
	if uid != comment.UserID {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	_, err = comment.DeleteAComment(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	response := map[string]interface{}{
		"token": "Comment has been deleted",
	}
	responses.JSON(w, http.StatusOK, response)
}

// GetComments godoc
// @Summary Get details of all comments
// @Description Get details of all comments
// @Tags comments
// @Param postId query int false "post's id comment belongs to"
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Comment
// @Router /api/comments/{postId} [get]
func (server *Server) GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["postId"], 10, 64)

	// check post exists
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	comment := models.Comment{}
	comments, err := comment.GetComments(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, comments)
}

// Update Comment godoc
// @Summary Update Comment By ID
// @Description Update details of a Comment by ID
// @Tags comments
// @Param id path int true "Comment ID"
// @Param Comment body models.Comment true "Update Request Body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Comment
// @Router /api/comments/{id} [put]
func (server *Server) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// validate token and get user Id from token
	uid, err := auth.ExtractTokenId(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	// validate comment exists
	origComment := models.Comment{}
	err = server.DB.Debug().Model(models.Comment{}).Where("id = ?", cid).Take(&origComment).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	if uid != origComment.UserID {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	// Read request data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// process request data
	commentUpdate := models.Comment{}
	err = json.Unmarshal(body, &commentUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	commentUpdate.Prepare()
	commentUpdate.ID = origComment.ID
	commentUpdate.UserID = origComment.UserID
	commentUpdate.PostID = origComment.PostID

	errorMessages := commentUpdate.Validate()
	if errorMessages != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	commentUpdated, err := commentUpdate.UpdateAComment(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, commentUpdated)
}
