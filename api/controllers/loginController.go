package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dmdinh22/go-blog/api/auth"
	"github.com/dmdinh22/go-blog/api/models"
	"github.com/dmdinh22/go-blog/api/responses"
	"github.com/dmdinh22/go-blog/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary Logs a user in
// @Description Authenticate credentials and logs user in
// @Tags login
// @Param login body models.Login true "auth token from login"
// @Accept  json
// @Produce  json
// @Success 200 {string} string "token"
// @Router /api/login [post]
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
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

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	bearerToken, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	response := map[string]interface{}{
		"token": bearerToken,
	}
	responses.JSON(w, http.StatusOK, response)
}

func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
