package x

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	// Create a temporary .env file for testing
	err := os.WriteFile(".env", []byte("TEST_VAR=test_value"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}
	defer os.Remove(".env")

	err = LoadEnv()
	if err != nil {
		t.Errorf("LoadEnv() returned an error: %v", err)
	}

	// Check if the environment variable was loaded
	if os.Getenv("TEST_VAR") != "test_value" {
		t.Errorf("LoadEnv() did not load the environment variable correctly")
	}
}

func TestLoadConfig(t *testing.T) {
	// Set test environment variables
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	}()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() returned an error: %v", err)
	}

	// Check if the config was loaded correctly
	if cfg.ServerPort != "8080" {
		t.Errorf("LoadConfig() ServerPort = %s, want %s", cfg.ServerPort, "8080")
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("LoadConfig() Database.Host = %s, want %s", cfg.Database.Host, "localhost")
	}
	if cfg.Database.Port != "5432" {
		t.Errorf("LoadConfig() Database.Port = %s, want %s", cfg.Database.Port, "5432")
	}

	if cfg.Database.Password != "testpass" {
		t.Errorf("LoadConfig() Database.Password = %s, want %s", cfg.Database.Password, "testpass")
	}

}

func TestConvertMapToJson(t *testing.T) {
	input := map[string]string{"key1": "value1", "key2": "value2"}
	expected := `{"key1":"value1","key2":"value2"}`

	result, err := ConvertMapToJson(input)
	if err != nil {
		t.Errorf("ConvertMapToJson() returned an error: %v", err)
	}
	if result != expected {
		t.Errorf("ConvertMapToJson() = %v, want %v", result, expected)
	}
}

func TestConvertJsonToMap(t *testing.T) {
	input := `{"key1":"value1","key2":"value2"}`
	expected := map[string]string{"key1": "value1", "key2": "value2"}

	result, err := ConvertJsonToMap(input)
	if err != nil {
		t.Errorf("ConvertJsonToMap() returned an error: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ConvertJsonToMap() = %v, want %v", result, expected)
	}
}

func TestGetStatusString(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		expected string
	}{
		{"Queued", 0, "QUEUED"},
		{"Running", 1, "RUNNING"},
		{"Failed", 2, "FAILED"},
		{"Succeeded", 3, "SUCCEEDED"},
		{"Unknown", 4, "UNKNOWN"},
		{"Negative", -1, "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStatusString(tt.status)
			if result != tt.expected {
				t.Errorf("GetStatusString(%d) = %s, want %s", tt.status, result, tt.expected)
			}
		})
	}
}

func TestCreateClient(t *testing.T) {
	// Test with default server URL
	t.Run("Default Server URL", func(t *testing.T) {
		os.Unsetenv("SERVER_ENDPOINT")
		client, err := CreateClient("http://task:8080")
		if err != nil {
			t.Fatalf("CreateClient() returned an error: %v", err)
		}
		if client == nil {
			t.Error("CreateClient() returned nil client")
		}
	})

	// Test with custom server URL
	t.Run("Custom Server URL", func(t *testing.T) {
		os.Setenv("SERVER_ENDPOINT", "http://testserver:9090")
		defer os.Unsetenv("SERVER_ENDPOINT")

		client, err := CreateClient("http://testserver:9090")
		if err != nil {
			t.Fatalf("CreateClient() returned an error: %v", err)
		}
		if client == nil {
			t.Error("CreateClient() returned nil client")
		}
		// Note: We can't easily check the actual URL used by the client,
		// but we can verify that a client was created without errors.
	})
}
