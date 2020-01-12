package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/dmdinh22/go-blog/api/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreatePost(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	token, err := server.SignIn(user.Email, "p@$$w0rd") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		title        string
		content      string
		authorId     uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"title":"The title", "content": "the content", "authorId": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			title:        "The title",
			content:      "the content",
			authorId:     user.ID,
			errorMessage: "",
		},
		{
			inputJSON:    `{"title":"The title", "content": "the content", "authorId": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title has already been used.",
		},
		{
			// When no token is passed
			inputJSON:    `{"title":"When no token is passed", "content": "the content", "authorId": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"title":"When incorrect token is passed", "content": "the content", "authorId": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"title": "", "content": "The content", "authorId": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			inputJSON:    `{"title": "This is a title", "content": "", "authorId": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			inputJSON:    `{"title": "This is an awesome title", "content": "the content"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"title": "This is an awesome title", "content": "the content", "authorId": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/posts", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreatePost)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			fmt.Print((responseMap))
			assert.Equal(t, responseMap["Title"], v.title)
			assert.Equal(t, responseMap["Content"], v.content)
			assert.Equal(t, responseMap["AuthorID"], float64(v.authorId)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetPosts(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetPosts)
	handler.ServeHTTP(rr, req)

	var posts []models.Post
	err = json.Unmarshal([]byte(rr.Body.String()), &posts)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(posts), 2)
}

func TestGetPostByID(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatal(err)
	}

	postSample := []struct {
		id           string
		statusCode   int
		title        string
		content      string
		authorId     uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(post.ID)),
			statusCode: 200,
			title:      post.Title,
			content:    post.Content,
			authorId:   post.AuthorID,
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}

	for _, v := range postSample {
		req, err := http.NewRequest("GET", "/posts", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetPost)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})

		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, post.Title, responseMap["Title"])
			assert.Equal(t, post.Content, responseMap["Content"])
			assert.Equal(t, float64(post.AuthorID), responseMap["AuthorID"]) //the response author id is float64
		}
	}
}
