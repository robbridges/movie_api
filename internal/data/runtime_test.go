package data

import (
	"encoding/json"
	"testing"
)

func TestRuntime_MarshalJSON(t *testing.T) {
	runtime := Runtime(90)

	jsonData, err := json.Marshal(runtime)
	if err != nil {
		t.Errorf("Failed to marshal Runtime: %v", err)
	}

	expectedJSON := `"90 mins"`
	if string(jsonData) != expectedJSON {
		t.Errorf("Unexpected JSON output. Expected: %s, Got: %s", expectedJSON, string(jsonData))
	}
}

func TestRuntime_UnmarshalJSON(t *testing.T) {
	jsonData := []byte(`"120 mins"`)

	var runtime Runtime
	err := json.Unmarshal(jsonData, &runtime)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	expectedRuntime := Runtime(120)
	if runtime != expectedRuntime {
		t.Errorf("Unexpected unmarshaled Runtime value. Expected: %d, Got: %d", expectedRuntime, runtime)
	}
}
