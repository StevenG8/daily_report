package collector

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"daily_report/internal/config"
)

func TestNewGitCollector(t *testing.T) {
	cfg := config.GitConfig{
		AuthorEmail: "test@example.com",
		Repos:       []string{},
		RepoDirs:    []string{},
	}

	collector := NewGitCollector(cfg)

	if collector == nil {
		t.Fatal("NewGitCollector returned nil")
	}

	if collector.Name() != "git" {
		t.Errorf("Expected Name() to return 'git', got '%s'", collector.Name())
	}
}

func TestGitCollector_Name(t *testing.T) {
	cfg := config.GitConfig{
		AuthorEmail: "test@example.com",
	}
	collector := NewGitCollector(cfg)

	name := collector.Name()
	if name != "git" {
		t.Errorf("Expected 'git', got '%s'", name)
	}
}

func TestGitCollector_GetRepositories(t *testing.T) {
	// Create temporary directory with git repos
	tempDir := t.TempDir()

	// Create fake git repos
	repo1 := filepath.Join(tempDir, "repo1")
	repo2 := filepath.Join(tempDir, "repo2")

	os.MkdirAll(filepath.Join(repo1, ".git"), 0755)
	os.MkdirAll(filepath.Join(repo2, ".git"), 0755)

	cfg := config.GitConfig{
		AuthorEmail: "test@example.com",
		RepoDirs:    []string{tempDir},
	}

	collector := NewGitCollector(cfg)
	repos, err := collector.getRepositories()

	if err != nil {
		t.Fatalf("getRepositories failed: %v", err)
	}

	if len(repos) != 2 {
		t.Errorf("Expected 2 repos, got %d", len(repos))
	}
}

func TestGitCollector_Collect_Empty(t *testing.T) {
	cfg := config.GitConfig{
		AuthorEmail: "test@example.com",
		Repos:       []string{},
		RepoDirs:    []string{},
	}

	collector := NewGitCollector(cfg)
	ctx := context.Background()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	items, err := collector.Collect(ctx, start, end)

	if err == nil {
		t.Error("Expected error for empty repos, got nil")
	}

	if items != nil {
		t.Errorf("Expected nil items, got %v", items)
	}
}

func TestGitCollector_ParseCommits(t *testing.T) {
	collector := &GitCollector{}

	output := `abc123|John Doe|john@example.com|2026-02-11 14:30:00 +0800|feat: add new feature
def456|Jane Smith|jane@example.com|2026-02-11 15:45:00 +0800|fix: fix bug`

	items, err := collector.parseCommits(output, "/path/to/repo")

	if err != nil {
		t.Fatalf("parseCommits failed: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	if items[0].Title != "feat: add new feature" {
		t.Errorf("Expected title 'feat: add new feature', got '%s'", items[0].Title)
	}

	if items[0].Type != "git" {
		t.Errorf("Expected type 'git', got '%s'", items[0].Type)
	}
}

func TestGitCollector_ParseCommits_Empty(t *testing.T) {
	collector := &GitCollector{}

	items, err := collector.parseCommits("", "/path/to/repo")

	if err != nil {
		t.Fatalf("parseCommits failed with empty input: %v", err)
	}

	if len(items) != 0 {
		t.Errorf("Expected 0 items for empty input, got %d", len(items))
	}
}
