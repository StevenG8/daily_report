package models

import "time"

// Item represents a single work output item from any data source
type Item struct {
	Type     string                 `json:"type"`               // "git", "meeting", "jira", "confluence"
	Title    string                 `json:"title"`              // Title or subject
	Time     time.Time              `json:"time"`               // Timestamp
	Link     string                 `json:"link"`               // URL to the item
	Content  string                 `json:"content,omitempty"`  // Detailed content (optional)
	Metadata map[string]interface{} `json:"metadata,omitempty"` // Extended fields
}

// ItemsByTime implements sort.Interface for sorting Items by time (descending)
type ItemsByTime []Item

func (it ItemsByTime) Len() int           { return len(it) }
func (it ItemsByTime) Swap(i, j int)      { it[i], it[j] = it[j], it[i] }
func (it ItemsByTime) Less(i, j int) bool { return it[i].Time.After(it[j].Time) }

// ReportData contains all collected items grouped by type
type ReportData struct {
	Date         time.Time               `json:"date"`
	StartTime    time.Time               `json:"start_time"`
	EndTime      time.Time               `json:"end_time"`
	Items        []Item                  `json:"items"`
	ItemsByType  map[string][]Item       `json:"items_by_type"`
	Stats        map[string]int          `json:"stats"`
	SourceStatus map[string]SourceStatus `json:"source_status"`
}

// SourceStatus represents the collection status of a data source
type SourceStatus struct {
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
