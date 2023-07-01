package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogError(t *testing.T) {
	// Create a mock HTTP request
	req, _ := http.NewRequest("GET", "/some-path", nil)

	// Create an instance of the application
	app := &application{}

	// Initialize a logger for the test
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.LstdFlags)
	app.logger = logger

	// Define a sample error
	err := errors.New("sample error")

	// Call the logError function
	app.logError(req, err)

	// Check the log output
	logOutput := buf.String()
	if !strings.Contains(logOutput, "sample error") {
		t.Error("Expected error message not found in the log output")
	}
}
func TestErrorResponse(t *testing.T) {
	// Create a mock HTTP response recorder
	recorder := httptest.NewRecorder()

	// Create a mock HTTP request
	req, _ := http.NewRequest("GET", "/some-path", nil)

	// Define the sample error message
	message := "Sample error message"

	// Call the errorResponse function
	app.errorResponse(recorder, req, http.StatusBadRequest, message)

	// Check the response status code
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	// Decode the response body
	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to decode response JSON: %v", err)
	}

	// Check the response data
	if errorMsg, ok := response["error"].(string); !ok || errorMsg != message {
		t.Errorf("Unexpected error message. Expected: %s, Got: %s", message, errorMsg)
	}
}

func TestServerErrorResponse(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/some-path", nil)

	var buf bytes.Buffer
	logger := log.New(&buf, "", log.LstdFlags)

	app := &application{
		logger: logger,
	}

	err := errors.New("sample error")

	app.serverErrorResponse(recorder, req, err)

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, recorder.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to decode response JSON: %v", err)
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "sample error") {
		t.Error("Expected error message not found in the log output")
	}
}
