package collector

import (
	"context"
	"time"

	"daily_report/pkg/models"
)

// Collector defines the interface for data collection from various sources
type Collector interface {
	// Name returns the name of the collector
	Name() string

	// Collect gathers items from the data source within the specified time range
	Collect(ctx context.Context, start, end time.Time) ([]models.Item, error)
}

// MultiCollector combines multiple collectors
type MultiCollector struct {
	collectors []Collector
}

// NewMultiCollector creates a new MultiCollector
func NewMultiCollector(collectors ...Collector) *MultiCollector {
	return &MultiCollector{
		collectors: collectors,
	}
}

// CollectAll collects items from all collectors
func (mc *MultiCollector) CollectAll(ctx context.Context, start, end time.Time) (map[string][]models.Item, map[string]models.SourceStatus) {
	result := make(map[string][]models.Item)
	status := make(map[string]models.SourceStatus)

	for _, c := range mc.collectors {
		items, err := c.Collect(ctx, start, end)
		if err != nil {
			status[c.Name()] = models.SourceStatus{
				Name:    c.Name(),
				Success: false,
				Error:   err.Error(),
			}
			continue
		}
		result[c.Name()] = items
		status[c.Name()] = models.SourceStatus{
			Name:    c.Name(),
			Success: true,
		}
	}

	return result, status
}
