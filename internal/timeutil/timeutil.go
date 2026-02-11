package timeutil

import (
	"strings"
	"time"
)

// GetDayRange returns the start and end time of the day for the given timezone
func GetDayRange(t time.Time, timezone string) (time.Time, time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Convert to local timezone
	localTime := t.In(loc)

	// Start of day: 00:00:00
	startOfDay := time.Date(
		localTime.Year(),
		localTime.Month(),
		localTime.Day(),
		0, 0, 0, 0,
		loc,
	)

	// End of day: 23:59:59
	endOfDay := time.Date(
		localTime.Year(),
		localTime.Month(),
		localTime.Day(),
		23, 59, 59, 999999999,
		loc,
	)

	return startOfDay, endOfDay, nil
}

// GetTodayRange returns the start and end time of today
func GetTodayRange(timezone string) (time.Time, time.Time, error) {
	return GetDayRange(time.Now(), timezone)
}

// ParseTimeRange parses custom time range strings
// Supported formats:
// - "today": today
// - "yesterday": yesterday
// - "2024-02-11": specific day
// - "2024-02-11,2024-02-12": custom range
func ParseTimeRange(input, timezone string) (time.Time, time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	switch input {
	case "today":
		return GetDayRange(time.Now(), timezone)
	case "yesterday":
		return GetDayRange(time.Now().AddDate(0, 0, -1), timezone)
	default:
		// Try to parse as single date or range
		if strings.Contains(input, ",") {
			parts := strings.Split(input, ",")
			if len(parts) == 2 {
				start, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(parts[0]), loc)
				if err != nil {
					return time.Time{}, time.Time{}, err
				}
				end, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(parts[1]), loc)
				if err != nil {
					return time.Time{}, time.Time{}, err
				}
				return start, time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, loc), nil
			}
		}

		// Try to parse as single date
		date, err := time.ParseInLocation("2006-01-02", input, loc)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return GetDayRange(date, timezone)
	}
}
