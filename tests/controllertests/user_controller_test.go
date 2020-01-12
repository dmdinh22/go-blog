package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON    string
		statusCode   int
		username     string
		email        string
		errorMessage string
	}{
		{
			inputJSON:    `{"username": "Pet", "email": "pet@gmail.com", "password": "p@$$w0rd"}`,
			statusCode:   201,
			username:     "Pet",
			email:        "pet@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"username": "Frank", "email": "pet@gmail.com", "password": "p@$$w0rd"}`,
			statusCode:   500,
			errorMessage: "Email has already been used.",
		},
		{
			inputJSON:    `{"username": "Pet", "email": "grand@gmail.com", "password": "p@$$w0rd"}`,
			statusCode:   500,
			errorMessage: "Username has already been used.",
		},
		{
			inputJSON:    `{"username": "Kan", "email": "kangmail.com", "password": "p@$$w0rd"}`,
			statusCode:   422,
			errorMessage: "Invalid Email.",
		},
		{
			inputJSON:    `{"username":  "", "email": "kan@gmail.com", "password": "p@$$w0rd"}`,
			statusCode:   422,
			errorMessage: "Required username.",
		},
		{
			inputJSON:    `{"username":  "Kan", "email": "", "password": "p@$$w0rd"}`,
			statusCode:   422,
			errorMessage: "Required Email.",
		},
		{
			inputJSON:    `{"username":  "Kan", "email": "kan@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "Required Password.",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})

		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["Username"], v.username)
			assert.Equal(t, responseMap["Email"], v.email)
		}

		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
