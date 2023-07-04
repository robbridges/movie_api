package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"movie_api/internal/validator"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var app = &application{}

func TestReadIDParam(t *testing.T) {

	t.Run("Valid ID", func(t *testing.T) {
		// Create a mock HTTP request with a valid ID parameter
		req, _ := http.NewRequest("GET", "/users/123", nil)

		// Create a new httprouter.Params instance with the mock ID parameter
		params := httprouter.Params{httprouter.Param{Key: "id", Value: "123"}}

		// Set the Params field of the request's context to the mock Params instance
		req = req.WithContext(addParamsToContext(req.Context(), params))

		// Call the readIDParam function
		id, err := app.readIDParam(req)

		// Assert that the returned ID is correct and the error is nil
		if id != 123 || err != nil {
			t.Errorf("Expected ID: 123, Got ID: %d, Error: %v", id, err)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Create a mock HTTP request with an invalid ID parameter
		req, _ := http.NewRequest("GET", "/users/abc", nil)

		// Set the Params field of the request's context to an empty Params instance
		req = req.WithContext(addParamsToContext(req.Context(), httprouter.Params{}))

		// Call the readIDParam function
		id, err := app.readIDParam(req)

		// Assert that the returned ID is 0 and the error is "invalid id parameter"
		if id != 0 || err == nil || err.Error() != "invalid id parameter" {
			t.Errorf("Expected ID: 0, Got ID: %d, Expected Error: 'invalid id parameter', Got Error: %v", id, err)
		}
	})
}

// Helper function to add httprouter.Params to the request context
func addParamsToContext(ctx context.Context, params httprouter.Params) context.Context {
	return context.WithValue(ctx, httprouter.ParamsKey, params)
}

func TestWriteJSON(t *testing.T) {

	recorder := httptest.NewRecorder()

	data := map[string]interface{}{
		"message": "Success",
		"data":    "Some data",
	}

	expectedJSON, _ := json.MarshalIndent(data, "", "\t")
	expectedJSON = append(expectedJSON, '\n')

	err := app.writeJson(recorder, http.StatusOK, data, nil)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	expectedHeaders := http.Header{"Content-Type": []string{"application/json"}}
	if !reflect.DeepEqual(recorder.Header(), expectedHeaders) {
		t.Errorf("Unexpected response headers: %v", recorder.Header())
	}

	if !reflect.DeepEqual(recorder.Body.Bytes(), expectedJSON) {
		t.Errorf("Unexpected response body:\nExpected: %s\nGot: %s", expectedJSON, recorder.Body.Bytes())
	}
}

func TestReadJSON(t *testing.T) {
	// Create a sample payload
	payload := []byte(`{"name":"John","age":30}`)

	// Create a mock HTTP request with the payload
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	// Create a mock HTTP response recorder
	recorder := httptest.NewRecorder()

	// Call the readJSON function
	var dst map[string]interface{}
	err := app.readJSON(recorder, req, &dst)

	// Check for any errors
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the expected values in the decoded payload
	expectedName := "John"
	expectedAge := 30

	if name, ok := dst["name"].(string); !ok || name != expectedName {
		t.Errorf("Unexpected name value. Expected: %s, Got: %s", expectedName, name)
	}

	if age, ok := dst["age"].(float64); !ok || int(age) != expectedAge {
		t.Errorf("Unexpected age value. Expected: %d, Got: %f", expectedAge, age)
	}
}

func TestReadString(t *testing.T) {
	app := &application{}

	// Test case 1: Value exists in query string
	qs := url.Values{"key": []string{"value"}}
	expected := "value"
	result := app.readString(qs, "key", "defaultValue")
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}

	// Test case 2: Value does not exist in query string, return default value
	qs = url.Values{}
	expected = "defaultValue"
	result = app.readString(qs, "key", "defaultValue")
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}
}

func TestReadCSV(t *testing.T) {
	app := &application{}

	// Test case 1: CSV value exists in query string
	qs := url.Values{"key": []string{"value1,value2,value3"}}
	expected := []string{"value1", "value2", "value3"}
	result := app.readCSV(qs, "key", []string{})
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case 2: CSV value does not exist in query string, return default value
	qs = url.Values{}
	defaultValue := []string{"default1", "default2"}
	expected = defaultValue
	result = app.readCSV(qs, "key", defaultValue)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestReadInt(t *testing.T) {
	app := &application{}
	v := validator.New()

	// Test case 1: Integer value exists in query string
	qs := url.Values{"key": []string{"42"}}
	expected := 42
	result := app.readInt(qs, "key", 0, v)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}

	// Test case 2: Integer value does not exist in query string, return default value
	qs = url.Values{}
	defaultValue := 123
	expected = defaultValue
	result = app.readInt(qs, "key", defaultValue, v)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}

}
