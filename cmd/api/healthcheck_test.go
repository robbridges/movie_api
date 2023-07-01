package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {

	// Create a mock HTTP request
	req, _ := http.NewRequest("GET", "/health", nil)

	// Create a mock HTTP response recorder
	recorder := httptest.NewRecorder()

	// Call the healthCheckHandler function
	app.healthCheckHandler(recorder, req)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	// Decode the response body
	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to decode response JSON: %v", err)
	}

	// Check the response data
	expectedStatus := "available"
	expectedVersion := "1.0.0" // Replace with your expected version value

	if status, ok := response["status"].(string); !ok || status != expectedStatus {
		t.Errorf("Unexpected status value. Expected: %s, Got: %s", expectedStatus, status)
	}

	if systemInfo, ok := response["system_info"].(map[string]interface{}); !ok {
		t.Error("Unexpected system_info value. Expected a map")
	} else {
		if version, ok := systemInfo["version"].(string); !ok || version != expectedVersion {
			t.Errorf("Unexpected version value. Expected: %s, Got: %s", expectedVersion, version)
		}
	}
}
