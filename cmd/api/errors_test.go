package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"movie_api/internal/jsonlog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLogError(t *testing.T) {
	// Create a mock HTTP request
	req, _ := http.NewRequest("GET", "/some-path", nil)

	app := &application{}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	app.logger = logger

	err := errors.New("sample error")

	app.logError(req, err)

	if !strings.Contains(err.Error(), "sample error") {
		t.Errorf("Expected error message not found")
	}

}
func TestErrorResponse(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/some-path", nil)

	message := "Sample error message"

	app.errorResponse(recorder, req, http.StatusBadRequest, message)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to decode response JSON: %v", err)
	}

	if errorMsg, ok := response["error"].(string); !ok || errorMsg != message {
		t.Errorf("Unexpected error message. Expected: %s, Got: %s", message, errorMsg)
	}
}

func TestServerErrorResponse(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/some-path", nil)

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

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

}

func TestErrorResponses(t *testing.T) {
	// Initialize the logger for the test
	var buf bytes.Buffer
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// Create an instance of the application with the necessary dependencies
	app := &application{
		logger: logger,
	}

	tests := []struct {
		name       string
		handler    func(w http.ResponseWriter, r *http.Request)
		statusCode int
		message    string
	}{
		{
			name:       "NotFoundResponse",
			handler:    app.notFoundResponse,
			statusCode: http.StatusNotFound,
			message:    "the requested resource could not be found",
		},
		{
			name:       "MethodNotAllowedResponse",
			handler:    app.methodNotAllowedResponse,
			statusCode: http.StatusMethodNotAllowed,
			message:    "The GET method is not supported for this resource",
		},
		{
			name: "BadRequestResponse",
			handler: func(w http.ResponseWriter, r *http.Request) {
				err := errors.New("bad request")
				app.badRequestResponse(w, r, err)
			},
			statusCode: http.StatusBadRequest,
			message:    "bad request",
		},

		{
			name:       "EditConflictResponse",
			handler:    app.editConflictResponse,
			statusCode: http.StatusConflict,
			message:    "unable to update the record due to an edit conflict, please try again",
		},
		{
			name:       "RateLimitError",
			handler:    app.rateLimitExceededResponse,
			statusCode: http.StatusTooManyRequests,
			message:    "rate limit exceeded",
		},
		{
			name:       "InvalidCredentialResponse",
			handler:    app.invalidCredentialResponse,
			statusCode: http.StatusUnauthorized,
			message:    "invalid credentials, please confirm and resubmit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock HTTP response recorder
			recorder := httptest.NewRecorder()

			// Create a mock HTTP request
			req, _ := http.NewRequest("GET", "/some-path", nil)

			// Call the handler function
			tt.handler(recorder, req)

			if recorder.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, recorder.Code)
			}

			var response map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to decode response JSON: %v", err)
			}

			if errorMsg, ok := response["error"].(string); !ok || errorMsg != tt.message {
				t.Errorf("Unexpected error message. Expected: %s, Got: %s", tt.message, errorMsg)
			}
		})
	}

	logOutput := buf.String()
	if logOutput != "" {
		t.Errorf("Unexpected log output: %s", logOutput)
	}
}
