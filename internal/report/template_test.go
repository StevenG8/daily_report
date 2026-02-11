package report

import (
	"testing"
	"time"

	"daily_report/pkg/models"
)

func TestNewGenerator(t *testing.T) {
	gen := NewGenerator()

	if gen == nil {
		t.Fatal("NewGenerator returned nil")
	}

	if gen.template != defaultTemplate {
		t.Error("Expected template to be set to default")
	}
}

func TestNewGeneratorWithTemplate(t *testing.T) {
	customTemplate := "# Custom Report\n\n{{date}}"
	gen := NewGeneratorWithTemplate(customTemplate)

	if gen == nil {
		t.Fatal("NewGeneratorWithTemplate returned nil")
	}

	if gen.template != customTemplate {
		t.Error("Template not set correctly")
	}
}

func TestNewGeneratorWithTemplatePath(t *testing.T) {
	gen := NewGeneratorWithTemplatePath("/path/to/template.md")

	if gen == nil {
		t.Fatal("NewGeneratorWithTemplatePath returned nil")
	}

	if gen.customTemplatePath != "/path/to/template.md" {
		t.Error("Template path not set correctly")
	}
}

func TestGenerator_CalculateStats(t *testing.T) {
	gen := NewGenerator()

	itemsByType := map[string][]models.Item{
		"git":     {{}, {}},
		"jira":    {{}},
		"meeting": []models.Item{},
	}

	stats := gen.calculateStats(itemsByType)

	if stats["git"] != 2 {
		t.Errorf("Expected git count 2, got %d", stats["git"])
	}

	if stats["jira"] != 1 {
		t.Errorf("Expected jira count 1, got %d", stats["jira"])
	}

	if stats["meeting"] != 0 {
		t.Errorf("Expected meeting count 0, got %d", stats["meeting"])
	}
}

func TestGenerator_GroupItemsByType(t *testing.T) {
	gen := NewGenerator()

	now := time.Now()
	items := []models.Item{
		{Type: "git", Time: now.Add(-1 * time.Hour)},
		{Type: "git", Time: now},
		{Type: "jira", Time: now.Add(-30 * time.Minute)},
	}

	grouped := gen.groupItemsByType(items)

	if len(grouped["git"]) != 2 {
		t.Errorf("Expected 2 git items, got %d", len(grouped["git"]))
	}

	if len(grouped["jira"]) != 1 {
		t.Errorf("Expected 1 jira item, got %d", len(grouped["jira"]))
	}

	// Check sorting (newest first)
	if grouped["git"][0].Time.Before(grouped["git"][1].Time) {
		t.Error("Items should be sorted by time descending")
	}
}

func TestGenerator_BuildSourceStatus(t *testing.T) {
	gen := NewGenerator()

	status := map[string]models.SourceStatus{
		"git": {
			Name:    "git",
			Success: true,
		},
		"jira": {
			Name:    "jira",
			Success: false,
			Error:   "API error",
		},
	}

	result := gen.buildSourceStatus(status)

	if result == "" {
		t.Error("Expected non-empty status string")
	}

	// Should contain both statuses
	if len(result) == 0 {
		t.Error("Status string should not be empty")
	}
}

func TestGenerator_Generate_Empty(t *testing.T) {
	gen := NewGenerator()

	reportData := &models.ReportData{
		Date:         time.Now(),
		StartTime:    time.Now().Add(-24 * time.Hour),
		EndTime:      time.Now(),
		Items:        []models.Item{},
		ItemsByType:  map[string][]models.Item{},
		Stats:        map[string]int{},
		SourceStatus: map[string]models.SourceStatus{},
	}

	result := gen.Generate(reportData)

	if result == "" {
		t.Error("Generate should return non-empty result")
	}
}

func TestGenerator_RenderCustomTemplate(t *testing.T) {
	gen := &Generator{}

	now := time.Now()
	itemsByType := map[string][]models.Item{
		"git": {
			{
				Type:  "git",
				Title: "Test commit",
				Time:  now,
				Metadata: map[string]interface{}{
					"repo":   "test-repo",
					"commit": "abc123456",
				},
			},
		},
	}
	stats := map[string]int{"git": 1}

	template := "# Report for {{date}}\n\nGit commits: {{git_count}}\n\n{{git_section}}"
	result := gen.renderCustomTemplate(&models.ReportData{}, itemsByType, stats, "", template)

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Should have replaced the placeholder
	if len(result) == len(template) {
		t.Error("Template placeholders should be replaced")
	}
}
