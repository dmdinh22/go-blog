package modeltests

import (
	"log"
	"testing"

	"github.com/dmdinh22/go-blog/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestGetAllUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	users, err := userInstance.GetAllUsers(server.DB)

	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, len(*users), 2)
}

func TestCreateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	newUser := models.User{
		ID:       1,
		Email:    "test@gmail.com",
		Username: "test",
		Password: "p@$$w0rd",
	}

	savedUser, err := newUser.CreateUser(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.Username, savedUser.Username)
}

func TestGetUserByID(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}

	foundUser, err := userInstance.GetUserById(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}

	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Username, user.Username)
}

func TestUpdateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}

	userUpdate := models.User{
		ID:       1,
		Username: "modiUpdate",
		Email:    "modiupdate@gmail.com",
		Password: "p@$$w0rd",
	}
	updatedUser, err := userUpdate.UpdateUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedUser.ID, userUpdate.ID)
	assert.Equal(t, updatedUser.Email, userUpdate.Email)
	assert.Equal(t, updatedUser.Username, userUpdate.Username)
}

func TestDeleteUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}

	isDeleted, err := userInstance.DeleteUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}

	//! one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}