package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"daily_report/internal/collector"
	"daily_report/internal/config"
	"daily_report/internal/report"
	"daily_report/internal/timeutil"
	"daily_report/pkg/models"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to config file")
	dateRange := flag.String("date", "today", "Date range: today, yesterday, or YYYY-MM-DD,YYYY-MM-DD")
	outputPath := flag.String("output", "", "Output file path (default: stdout)")
	mode := flag.String("mode", "template", "Report mode: template or llm")
	templatePath := flag.String("template", "", "Path to custom Markdown template file")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Daily Report Generator\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  daily_report [options]\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nExamples:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  daily_report                          # Generate today's report\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  daily_report --date yesterday        # Generate yesterday's report\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  daily_report --output report.md      # Save to file\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  daily_report --template custom.tmpl  # Use custom template\n")
	}
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Override mode from command line
	if *mode != "" {
		cfg.Report.Mode = *mode
	}

	// Override template path from command line
	if *templatePath != "" {
		cfg.Report.TemplatePath = *templatePath
	}

	// Parse time range
	start, end, err := timeutil.ParseTimeRange(*dateRange, cfg.Time.Timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing time range: %v\n", err)
		os.Exit(1)
	}

	// Create collectors
	gitCollector := collector.NewGitCollector(cfg.Git)
	multiCollector := collector.NewMultiCollector(gitCollector)

	// Collect data
	ctx := context.Background()
	itemsByType, sourceStatus := multiCollector.CollectAll(ctx, start, end)

	// Flatten items
	var allItems []models.Item
	for _, items := range itemsByType {
		allItems = append(allItems, items...)
	}

	// Prepare report data
	reportData := &models.ReportData{
		Date:         start,
		StartTime:    start,
		EndTime:      end,
		Items:        allItems,
		ItemsByType:  itemsByType,
		Stats:        make(map[string]int),
		SourceStatus: sourceStatus,
	}

	// Calculate stats
	for typ, items := range itemsByType {
		reportData.Stats[typ] = len(items)
	}

	// Generate report
	var generator *report.Generator
	if cfg.Report.TemplatePath != "" {
		generator = report.NewGeneratorWithTemplatePath(cfg.Report.TemplatePath)
	} else {
		generator = report.NewGenerator()
	}
	markdown := generator.Generate(reportData)

	// Output report
	if *outputPath != "" {
		if err := os.WriteFile(*outputPath, []byte(markdown), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Report generated: %s\n", *outputPath)
	} else {
		fmt.Print(markdown)
	}
}
