package controllers

import (
	"net/http"

	"github.com/dmdinh22/go-blog/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, r*http.Request)
}
