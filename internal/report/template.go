package report

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"daily_report/pkg/models"
)

// Generator generates daily reports from collected data
type Generator struct {
	template           string
	customTemplatePath string
}

// NewGenerator creates a new report generator with default template
func NewGenerator() *Generator {
	return &Generator{
		template: defaultTemplate,
	}
}

// NewGeneratorWithTemplate creates a new report generator with custom template
func NewGeneratorWithTemplate(template string) *Generator {
	return &Generator{
		template: template,
	}
}

// NewGeneratorWithTemplatePath creates a new report generator that loads template from file
func NewGeneratorWithTemplatePath(templatePath string) *Generator {
	return &Generator{
		customTemplatePath: templatePath,
	}
}

// Generate generates a markdown report from the collected data
func (g *Generator) Generate(data *models.ReportData) string {
	// Load template
	template := g.getTemplate()

	// Group items by type
	itemsByType := g.groupItemsByType(data.Items)

	// Calculate stats
	stats := g.calculateStats(itemsByType)

	// Build source status
	sourceStatus := g.buildSourceStatus(data.SourceStatus)

	// Render template
	return g.renderTemplate(data, itemsByType, stats, sourceStatus, template)
}

// getTemplate returns the template to use (from file or default)
func (g *Generator) getTemplate() string {
	if g.customTemplatePath != "" {
		data, err := os.ReadFile(g.customTemplatePath)
		if err != nil {
			// Fall back to default template if file read fails
			return defaultTemplate
		}
		return string(data)
	}
	return g.template
}

// groupItemsByType groups items by their type and sorts by time
func (g *Generator) groupItemsByType(items []models.Item) map[string][]models.Item {
	result := make(map[string][]models.Item)

	for _, item := range items {
		result[item.Type] = append(result[item.Type], item)
	}

	// Sort items by time (descending)
	for typ := range result {
		sort.Sort(models.ItemsByTime(result[typ]))
	}

	return result
}

// calculateStats calculates statistics for each item type
func (g *Generator) calculateStats(itemsByType map[string][]models.Item) map[string]int {
	stats := make(map[string]int)

	for typ, items := range itemsByType {
		stats[typ] = len(items)
	}

	return stats
}

// buildSourceStatus builds a formatted source status string
func (g *Generator) buildSourceStatus(status map[string]models.SourceStatus) string {
	var parts []string

	// Sort by source name for consistent output
	sources := make([]string, 0, len(status))
	for name := range status {
		sources = append(sources, name)
	}
	sort.Strings(sources)

	for _, name := range sources {
		s := status[name]
		if s.Success {
			parts = append(parts, fmt.Sprintf("✅ %s", s.Name))
		} else {
			parts = append(parts, fmt.Sprintf("❌ %s (%s)", s.Name, s.Error))
		}
	}

	return strings.Join(parts, " | ")
}

// renderTemplate renders the report using the template
func (g *Generator) renderTemplate(data *models.ReportData, itemsByType map[string][]models.Item, stats map[string]int, sourceStatus string, template string) string {
	// For custom templates, use simple placeholder replacement
	if g.customTemplatePath != "" {
		return g.renderCustomTemplate(data, itemsByType, stats, sourceStatus, template)
	}

	// Use default template renderer
	var sb strings.Builder

	// Title and date
	sb.WriteString(fmt.Sprintf("# 日报 - %s\n\n", data.Date.Format("2006年1月2日")))

	// Summary stats
	sb.WriteString("## 📊 汇总统计\n\n")
	sb.WriteString(fmt.Sprintf("- Git 提交: %d 次\n", stats["git"]))
	sb.WriteString(fmt.Sprintf("- 会议: %d 场\n", stats["meeting"]))
	sb.WriteString(fmt.Sprintf("- Jira 任务: %d 个\n", stats["jira"]))
	sb.WriteString(fmt.Sprintf("- Confluence 文档: %d 篇\n\n", stats["confluence"]))

	// Git commits
	if items, ok := itemsByType["git"]; ok && len(items) > 0 {
		sb.WriteString(g.renderGitItems(items))
	}

	// Meetings
	if items, ok := itemsByType["meeting"]; ok && len(items) > 0 {
		sb.WriteString(g.renderMeetingItems(items))
	}

	// Jira tasks
	if items, ok := itemsByType["jira"]; ok && len(items) > 0 {
		sb.WriteString(g.renderJiraItems(items))
	}

	// Confluence docs
	if items, ok := itemsByType["confluence"]; ok && len(items) > 0 {
		sb.WriteString(g.renderConfluenceItems(items))
	}

	// Footer
	sb.WriteString(fmt.Sprintf("\n---\n\n生成于: %s\n数据源状态: %s\n",
		time.Now().Format("2006-01-02 15:04:05"), sourceStatus))

	return sb.String()
}

// renderCustomTemplate renders using custom template with placeholders
func (g *Generator) renderCustomTemplate(data *models.ReportData, itemsByType map[string][]models.Item, stats map[string]int, sourceStatus string, template string) string {
	result := template

	// Replace basic placeholders
	result = strings.ReplaceAll(result, "{{date}}", data.Date.Format("2006年1月2日"))
	result = strings.ReplaceAll(result, "{{date_en}}", data.Date.Format("2006-01-02"))
	result = strings.ReplaceAll(result, "{{generate_time}}", time.Now().Format("2006-01-02 15:04:05"))
	result = strings.ReplaceAll(result, "{{source_status}}", sourceStatus)

	// Replace stats
	result = strings.ReplaceAll(result, "{{git_count}}", fmt.Sprintf("%d", stats["git"]))
	result = strings.ReplaceAll(result, "{{meeting_count}}", fmt.Sprintf("%d", stats["meeting"]))
	result = strings.ReplaceAll(result, "{{jira_count}}", fmt.Sprintf("%d", stats["jira"]))
	result = strings.ReplaceAll(result, "{{confluence_count}}", fmt.Sprintf("%d", stats["confluence"]))

	// Replace sections
	result = strings.ReplaceAll(result, "{{git_section}}", g.renderGitItems(itemsByType["git"]))
	result = strings.ReplaceAll(result, "{{meeting_section}}", g.renderMeetingItems(itemsByType["meeting"]))
	result = strings.ReplaceAll(result, "{{jira_section}}", g.renderJiraItems(itemsByType["jira"]))
	result = strings.ReplaceAll(result, "{{confluence_section}}", g.renderConfluenceItems(itemsByType["confluence"]))

	return result
}

// renderGitItems renders Git commit items
func (g *Generator) renderGitItems(items []models.Item) string {
	var sb strings.Builder

	sb.WriteString("## 💻 代码提交\n\n")

	// Group by repo
	byRepo := make(map[string][]models.Item)
	for _, item := range items {
		repo := item.Metadata["repo"].(string)
		byRepo[repo] = append(byRepo[repo], item)
	}

	for repo, commits := range byRepo {
		sb.WriteString(fmt.Sprintf("### %s\n\n", repo))
		for _, commit := range commits {
			sb.WriteString(fmt.Sprintf("- %s (%s)\n",
				commit.Title,
				commit.Time.Format("15:04")))
			if commitLink, ok := commit.Metadata["commit"].(string); ok {
				sb.WriteString(fmt.Sprintf("  commit: %s\n", commitLink[:7]))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderMeetingItems renders meeting items
func (g *Generator) renderMeetingItems(items []models.Item) string {
	var sb strings.Builder

	sb.WriteString("## 📅 会议\n\n")

	for _, item := range items {
		sb.WriteString(fmt.Sprintf("### %s - %s\n",
			item.Time.Format("15:04"),
			item.Title))
		sb.WriteString(fmt.Sprintf("- 参会者: %s\n", item.Content))
		if item.Link != "" {
			sb.WriteString(fmt.Sprintf("- 链接: %s\n", item.Link))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderJiraItems renders Jira task items
func (g *Generator) renderJiraItems(items []models.Item) string {
	var sb strings.Builder

	sb.WriteString("## 🎯 Jira 任务\n\n")

	for _, item := range items {
		sb.WriteString(fmt.Sprintf("### %s\n", item.Title))
		if status, ok := item.Metadata["status"].(string); ok {
			sb.WriteString(fmt.Sprintf("- 状态: %s\n", status))
		}
		sb.WriteString(fmt.Sprintf("- 更新时间: %s\n", item.Time.Format("15:04")))
		if item.Link != "" {
			sb.WriteString(fmt.Sprintf("- 链接: %s\n", item.Link))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderConfluenceItems renders Confluence document items
func (g *Generator) renderConfluenceItems(items []models.Item) string {
	var sb strings.Builder

	sb.WriteString("## 📝 Confluence 文档\n\n")

	for _, item := range items {
		sb.WriteString(fmt.Sprintf("### %s\n", item.Title))
		sb.WriteString(fmt.Sprintf("- 作者: %s\n", item.Content))
		sb.WriteString(fmt.Sprintf("- 更新时间: %s\n", item.Time.Format("15:04")))
		if item.Link != "" {
			sb.WriteString(fmt.Sprintf("- 链接: %s\n", item.Link))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// defaultTemplate is the default markdown template
const defaultTemplate = ""
