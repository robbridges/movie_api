package main

import (
	"movie_api/internal/data"
	"net/http"
	"testing"
)

func TestApplication_ContextUser(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
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
	})

	t.Run("Sad path, panic", func(t *testing.T) {
		app := application{}

		badRequest, _ := http.NewRequest(http.MethodGet, "/", nil)

		defer func() {
			if err := recover(); err == nil {
				t.Error("Expected a panic, but no panic occurred")
			}
		}()

		_ = app.contextGetUser(badRequest)

	})
}
