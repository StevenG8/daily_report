package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Git        GitConfig        `yaml:"git"`
	Meetings   MeetingsConfig   `yaml:"meetings"`
	Jira       JiraConfig       `yaml:"jira"`
	Confluence ConfluenceConfig `yaml:"confluence"`
	Report     ReportConfig     `yaml:"report"`
	Time       TimeConfig       `yaml:"time"`
}

// GitConfig contains Git collector configuration
type GitConfig struct {
	Author   string   `yaml:"author"`
	Repos    []string `yaml:"repos"`     // Specific repository paths
	RepoDirs []string `yaml:"repo_dirs"` // Directories to scan for git repos
}

// MeetingsConfig contains meeting collector configuration
type MeetingsConfig struct {
	Platform  string `yaml:"platform"` // feishu, dingtalk, wecom
	UserID    string `yaml:"user_id"`
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

// JiraConfig contains Jira collector configuration
type JiraConfig struct {
	Username   string `yaml:"username"`
	URL        string `yaml:"url"`
	APIToken   string `yaml:"api_token"`
	ProjectKey string `yaml:"project_key"`
}

// ConfluenceConfig contains Confluence collector configuration
type ConfluenceConfig struct {
	Username string `yaml:"username"`
	URL      string `yaml:"url"`
	APIToken string `yaml:"api_token"`
	SpaceKey string `yaml:"space_key"`
}

// ReportConfig contains report generation configuration
type ReportConfig struct {
	Mode         string    `yaml:"mode"`          // template, llm
	TemplatePath string    `yaml:"template_path"` // Path to custom template file
	LLM          LLMConfig `yaml:"llm"`
}

// LLMConfig contains LLM configuration for report generation
type LLMConfig struct {
	Provider     string `yaml:"provider"` // openai, anthropic, local
	Model        string `yaml:"model"`
	APIKey       string `yaml:"api_key"`
	SystemPrompt string `yaml:"system_prompt"`
}

// TimeConfig contains time configuration
type TimeConfig struct {
	Timezone string `yaml:"timezone"`
}

// Load loads configuration from a YAML file and expands environment variables
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables
	expanded := expandEnvVars(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set default values
	if cfg.Time.Timezone == "" {
		cfg.Time.Timezone = "Asia/Shanghai"
	}

	return &cfg, nil
}

// expandEnvVars replaces ${VAR_NAME} with environment variable values
func expandEnvVars(input string) string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		// Extract variable name from ${VAR_NAME}
		varName := strings.Trim(match, "${}")
		// Get environment variable value
		if val := os.Getenv(varName); val != "" {
			return val
		}
		// Return original if not found
		return match
	})
}
