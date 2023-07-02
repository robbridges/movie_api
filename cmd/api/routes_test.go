package main

import "testing"

func TestRoutes(t *testing.T) {

	app := &application{}

	router := app.routes()

	if router == nil {
		t.Error("Expected a non-nil router, got nil")
	}

	if router.NotFound == nil || router.MethodNotAllowed == nil {
		t.Error("Expected NotFound and MethodNotAllowed handlers to be set")
	}
	
}
