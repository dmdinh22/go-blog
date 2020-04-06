package controllers

import (
	"net/http"

	"github.com/dmdinh22/go-blog/api/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/dmdinh22/go-blog/docs"
)

func (s *Server) initializeRoutes() {
	// Home Route
	s.Router.HandleFunc("/api", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/api/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	// User routes
	s.Router.HandleFunc("/api/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/api/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/api/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/api/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/api/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Post routes
	s.Router.HandleFunc("/api/posts", middlewares.SetMiddlewareJSON(s.CreatePost)).Methods("POST")
	s.Router.HandleFunc("/api/posts", middlewares.SetMiddlewareJSON(s.GetPosts)).Methods("GET")
	s.Router.HandleFunc("/api/posts/{id}", middlewares.SetMiddlewareJSON(s.GetPost)).Methods("GET")
	s.Router.HandleFunc("/api/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")
	s.Router.HandleFunc("/api/posts/{id}", middlewares.SetMiddlewareAuthentication(s.DeletePost)).Methods("DELETE")

	//Comment routes
	s.Router.HandleFunc("/api/comments", middlewares.SetMiddlewareJSON(s.CreateComment)).Methods("POST")

	// Swagger
	s.Router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// redirect base route to swagger
	s.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/"+"swagger/", http.StatusSeeOther)
	})
}
