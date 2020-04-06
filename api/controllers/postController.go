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

// CreatePost godoc
// @Summary Creates a new post
// @Description Creates a new post for the user logged in
// @Tags posts
// @Param ID query int false "post's id "
// @Param Title query string true "post's title"
// @Param Content query string true "post's content"
// @Param AuthorID query string true "id of user that created this post"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Post
// @Router /api/posts [post]
func (server *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	post := models.Post{}
	err = json.Unmarshal(body, &post)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	post.Prepare()
	err = post.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenId(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	createdPost, err := post.CreatePost(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, createdPost.ID))

	responses.JSON(w, http.StatusCreated, createdPost)
}

// GetPosts godoc
// @Summary Get details of all posts
// @Description Get details of all posts
// @Tags posts
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Post
// @Router /api/posts [get]
func (server *Server) GetPosts(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}
	posts, err := post.GetAllPosts(server.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

// GetPost godoc
// @Summary Get post By ID
// @Description Get details of a post by ID
// @Tags posts
// @Param id path int true "post ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Post
// @Router /api/posts/{id} [get]
func (server *Server) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	post := models.Post{}

	postRetrieved, err := post.GetPostByID(server.DB, pid)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, postRetrieved)
}

// Update Post godoc
// @Summary Update Post By ID
// @Description Update details of a Post by ID
// @Tags posts
// @Param id path int true "Post ID"
// @Param Post body models.Post true "Update Request Body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Post
// @Router /api/posts/{id} [put]
func (server *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if  postId is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check if auth token valid & get the userId
	uid, err := auth.ExtractTokenId(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if post exists
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post not found"))
		return
	}

	// Check if authorized author
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != postUpdate.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	postUpdate.Prepare()
	//this is important to tell the model the post id to update, the other update field are set above
	postUpdate.ID = post.ID
	err = postUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedPost, err := postUpdate.UpdatePost(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, updatedPost)
}

// Delete Post godoc
// @Summary Delete Post By ID
// @Description Delete details of a Post by ID
// @Tags posts
// @Param id path int true "Post ID"
// @Accept  json
// @Produce  json
// @Success 200
// @Router /api/posts/{id} [delete]
func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check valid post id
	pid, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check if user authenticated?
	uid, err := auth.ExtractTokenId(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exists
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this post?
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	_, err = post.DeletePost(server.DB, pid, uid)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	response := map[string]interface{}{
		"token": "Post has been deleted",
	}
	responses.JSON(w, http.StatusOK, response)
}
