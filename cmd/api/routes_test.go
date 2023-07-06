package main

import "testing"

func TestRoutes(t *testing.T) {

	app := &application{}

	router := app.routes()

	if router == nil {
		t.Error("Expected a non-nil router, got nil")
	}

}
