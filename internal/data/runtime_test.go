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
	t.Run("Valid json happy path", func(t *testing.T) {
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
	})
}

func TestRuntime_UnmarshalJSONSadPaths(t *testing.T) {
	tests := []struct {
		Name     string
		jsonData []byte
	}{
		{"Invalid Json not numberical",
			[]byte(`"abc mins"`),
		},
		{
			"Invalid Json doesn't end in mins",
			[]byte(`"123 minsx"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var runtime Runtime
			err := json.Unmarshal(tt.jsonData, &runtime)
			if err == nil {
				t.Errorf("Expected error, but not none")
			}
		})
	}
}
