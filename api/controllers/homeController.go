package controllers

import (
	"net/http"

	"github.com/dmdinh22/go-blog/api/responses"
)

// Home godoc
// @Summary Main route to check API is running
// @Produce json
// @Tags home
// @Success 200
// @Router /api [get]
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Go ftw!")
}
