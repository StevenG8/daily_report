package collector

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"daily_report/internal/config"
	"daily_report/pkg/models"
)

// GitCollector collects Git commits from repositories
type GitCollector struct {
	cfg config.GitConfig
}

// NewGitCollector creates a new Git collector
func NewGitCollector(cfg config.GitConfig) *GitCollector {
	return &GitCollector{cfg: cfg}
}

// Name returns the name of the collector
func (g *GitCollector) Name() string {
	return "git"
}

// Collect gathers Git commits from configured repositories
func (g *GitCollector) Collect(ctx context.Context, start, end time.Time) ([]models.Item, error) {
	// Get all repositories
	repos, err := g.getRepositories()
	if err != nil {
		return nil, fmt.Errorf("failed to get repositories: %w", err)
	}

	if len(repos) == 0 {
		return nil, fmt.Errorf("no repositories found")
	}

	var allItems []models.Item

	for _, repo := range repos {
		items, err := g.collectFromRepo(ctx, repo, start, end)
		if err != nil {
			return nil, fmt.Errorf("failed to collect from repo %s: %w", repo, err)
		}
		allItems = append(allItems, items...)
	}

	return allItems, nil
}

// getRepositories returns all repositories from explicit repos and scanned directories
func (g *GitCollector) getRepositories() ([]string, error) {
	repos := make(map[string]bool)

	// Add explicit repos
	for _, repo := range g.cfg.Repos {
		repos[repo] = true
	}

	// Scan directories for git repos
	for _, dir := range g.cfg.RepoDirs {
		foundRepos, err := g.scanDirectory(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to scan directory %s: %w", dir, err)
		}
		for _, repo := range foundRepos {
			repos[repo] = true
		}
	}

	// Convert to slice
	result := make([]string, 0, len(repos))
	for repo := range repos {
		result = append(result, repo)
	}

	return result, nil
}

// scanDirectory scans a directory for git repositories
func (g *GitCollector) scanDirectory(dir string) ([]string, error) {
	var repos []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if this is a .git directory
		if info.IsDir() && info.Name() == ".git" {
			// The parent directory is the git repository
			repoPath := filepath.Dir(path)
			repos = append(repos, repoPath)
			// Skip traversing into this directory
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return repos, nil
}

// collectFromRepo collects commits from a single repository
func (g *GitCollector) collectFromRepo(ctx context.Context, repo string, start, end time.Time) ([]models.Item, error) {
	// Check if repo exists
	cmd := exec.CommandContext(ctx, "git", "-C", repo, "rev-parse", "--git-dir")
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("not a git repository: %s", repo)
	}

	// Get commits with author filter
	startStr := start.Format("2006-01-02 15:04:05")
	endStr := end.Format("2006-01-02 15:04:05")

	cmd = exec.CommandContext(ctx, "git", "-C", repo, "log",
		"--author="+g.cfg.Author,
		"--since="+startStr,
		"--until="+endStr,
		"--pretty=format:%H|%an|%ae|%ai|%s")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git log: %w", err)
	}

	return g.parseCommits(string(output), repo)
}

// parseCommits parses git log output into Items
func (g *GitCollector) parseCommits(output, repo string) ([]models.Item, error) {
	if output == "" {
		return []models.Item{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	items := make([]models.Item, 0, len(lines))

	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue
		}

		commitHash := parts[0]
		authorName := parts[1]
		authorEmail := parts[2]
		commitTime, err := time.Parse("2006-01-02 15:04:05 -0700", parts[3])
		if err != nil {
			continue
		}
		message := parts[4]

		// Get repo name from path
		repoName := repo
		if idx := strings.LastIndex(repo, "/"); idx != -1 {
			repoName = repo[idx+1:]
		}

		item := models.Item{
			Type:    "git",
			Title:   message,
			Time:    commitTime,
			Link:    fmt.Sprintf("%s/commit/%s", repo, commitHash),
			Content: fmt.Sprintf("%s <%s>", authorName, authorEmail),
			Metadata: map[string]interface{}{
				"repo":         repoName,
				"commit":       commitHash,
				"author":       authorName,
				"author_email": authorEmail,
			},
		}

		items = append(items, item)
	}

	return items, nil
}
