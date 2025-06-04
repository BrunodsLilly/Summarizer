package core

import (
	"os"
	"testing"

	gemini_api "github.com/BrunodsLilly/Summarizer/pkg/core/internal/utils"
)

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Fatal("NewApp() returned nil")
	}
}

func TestGetModelName(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "uses environment variable when set",
			envValue: "custom-model",
			expected: "custom-model",
		},
		{
			name:     "uses default when env var not set",
			envValue: "",
			expected: ModelName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("MODEL_NAME", tt.envValue)
				defer os.Unsetenv("MODEL_NAME")
			} else {
				os.Unsetenv("MODEL_NAME")
			}

			result := gemini_api.GetModelName()
			if result != tt.expected {
				t.Errorf("GetModelName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetAPIVersion(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "uses environment variable when set",
			envValue: "v2beta",
			expected: "v2beta",
		},
		{
			name:     "uses default when env var not set",
			envValue: "",
			expected: APIVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("API_VERSION", tt.envValue)
				defer os.Unsetenv("API_VERSION")
			} else {
				os.Unsetenv("API_VERSION")
			}

			result := gemini_api.GetAPIVersion()
			if result != tt.expected {
				t.Errorf("getAPIVersion() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSummarizeURL_EmptyURL(t *testing.T) {
	_, err := SummarizeURL("")
	if err == nil {
		t.Error("SummarizeURL(\"\") should return an error for empty URL")
	}
}

func TestSummarizeURL_InvalidURL(t *testing.T) {
	_, err := SummarizeURL("not-a-valid-url")
	if err == nil {
		t.Error("SummarizeURL() should return an error for invalid URL")
	}
}
