package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dmdinh22/go-blog/api/auth"
	"github.com/dmdinh22/go-blog/api/models"
	"github.com/dmdinh22/go-blog/api/responses"
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

// func (server *Server) DeleteComment(w http.ResponseWriter, r *http.Request) {
// 	commentID := c.Param("id")
// 	// Is a valid post id given to us?
// 	cid, err := strconv.ParseUint(commentID, 10, 64)
// 	if err != nil {
// 		errList["Invalid_request"] = "Invalid Request"
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status": http.StatusBadRequest,
// 			"error":  errList,
// 		})
// 		return
// 	}
// 	// Is this user authenticated?
// 	uid, err := auth.ExtractTokenId(c.Request)
// 	if err != nil {
// 		errList["Unauthorized"] = "Unauthorized"
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"status": http.StatusUnauthorized,
// 			"error":  errList,
// 		})
// 		return
// 	}
// 	// Check if the comment exist
// 	comment := models.Comment{}
// 	err = server.DB.Debug().Model(models.Comment{}).Where("id = ?", cid).Take(&comment).Error
// 	if err != nil {
// 		errList["No_post"] = "No Post Found"
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"status": http.StatusNotFound,
// 			"error":  errList,
// 		})
// 		return
// 	}
// 	// Is the authenticated user, the owner of this post?
// 	if uid != comment.UserID {
// 		errList["Unauthorized"] = "Unauthorized"
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"status": http.StatusUnauthorized,
// 			"error":  errList,
// 		})
// 		return
// 	}

// 	// If all the conditions are met, delete the post
// 	_, err = comment.DeleteAComment(server.DB)
// 	if err != nil {
// 		errList["Other_error"] = "Please try again later"
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"status": http.StatusNotFound,
// 			"error":  errList,
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":   http.StatusOK,
// 		"response": "Comment deleted",
// 	})
// }
