package main

import (
	"movie_api/internal/data"
	"net/http"
	"testing"
)

func TestApplication_ContextSetUser(t *testing.T) {
	app := application{}

	user := &data.User{
		ID: 1,
	}

	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	newRequest := app.contextSetUser(r, user)

	retrievedUser := app.contextGetUser(newRequest)
	if retrievedUser != user {
		t.Errorf("Users were supposed to match")
	}
}
