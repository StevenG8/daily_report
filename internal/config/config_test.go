package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandEnvVars(t *testing.T) {
	// Set test environment variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single variable",
			input:    "${TEST_VAR}",
			expected: "test_value",
		},
		{
			name:     "Multiple variables",
			input:    "${TEST_VAR} and ${TEST_VAR}",
			expected: "test_value and test_value",
		},
		{
			name:     "No variables",
			input:    "plain text",
			expected: "plain text",
		},
		{
			name:     "Mixed variables and text",
			input:    "prefix ${TEST_VAR} suffix",
			expected: "prefix test_value suffix",
		},
		{
			name:     "Undefined variable",
			input:    "${UNDEFINED_VAR}",
			expected: "${UNDEFINED_VAR}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandEnvVars(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestLoad_InvalidPath(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")

	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Write invalid YAML
	os.WriteFile(configPath, []byte("invalid: yaml: content:"), 0644)

	_, err := Load(configPath)

	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestLoad_ValidYAML(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	yamlContent := `
git:
  author_email: "test@example.com"
  repos: []
  repo_dirs: []

meetings:
  platform: "feishu"
  user_id: "test_user"

jira:
  username: "test_user"

confluence:
  username: "test_user"

report:
  mode: "template"

time:
  timezone: "Asia/Shanghai"
`

	os.WriteFile(configPath, []byte(yamlContent), 0644)

	cfg, err := Load(configPath)

	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Git.AuthorEmail != "test@example.com" {
		t.Errorf("Expected author_email 'test@example.com', got '%s'", cfg.Git.AuthorEmail)
	}

	if cfg.Time.Timezone != "Asia/Shanghai" {
		t.Errorf("Expected timezone 'Asia/Shanghai', got '%s'", cfg.Time.Timezone)
	}
}

func TestLoad_EnvVars(t *testing.T) {
	// Set test environment variable
	os.Setenv("TEST_API_TOKEN", "secret_token")
	defer os.Unsetenv("TEST_API_TOKEN")

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	yamlContent := `
git:
  author_email: "test@example.com"
  repos: []
  repo_dirs: []

jira:
  username: "test_user"
  api_token: "${TEST_API_TOKEN}"

report:
  mode: "template"

time:
  timezone: "Asia/Shanghai"
`

	os.WriteFile(configPath, []byte(yamlContent), 0644)

	cfg, err := Load(configPath)

	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Jira.APIToken != "secret_token" {
		t.Errorf("Expected api_token 'secret_token', got '%s'", cfg.Jira.APIToken)
	}
}

func TestLoad_DefaultTimezone(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	yamlContent := `
git:
  author_email: "test@example.com"
  repos: []
  repo_dirs: []

report:
  mode: "template"

time:
  timezone: ""
`

	os.WriteFile(configPath, []byte(yamlContent), 0644)

	cfg, err := Load(configPath)

	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Time.Timezone != "Asia/Shanghai" {
		t.Errorf("Expected default timezone 'Asia/Shanghai', got '%s'", cfg.Time.Timezone)
	}
}
