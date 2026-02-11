package timeutil

import (
	"testing"
	"time"
)

func TestGetDayRange(t *testing.T) {
	tests := []struct {
		name      string
		date      time.Time
		timezone  string
		wantErr   bool
		checkHour bool // Check if hours are 0 and 23
	}{
		{
			name:      "Valid date with Shanghai timezone",
			date:      time.Date(2026, 2, 11, 14, 30, 0, 0, time.UTC),
			timezone:  "Asia/Shanghai",
			wantErr:   false,
			checkHour: true,
		},
		{
			name:      "Valid date with UTC timezone",
			date:      time.Date(2026, 2, 11, 14, 30, 0, 0, time.UTC),
			timezone:  "UTC",
			wantErr:   false,
			checkHour: true,
		},
		{
			name:     "Invalid timezone",
			date:     time.Now(),
			timezone: "Invalid/Timezone",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := GetDayRange(tt.date, tt.timezone)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetDayRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkHour {
				// Check start is at 00:00:00
				if start.Hour() != 0 || start.Minute() != 0 || start.Second() != 0 {
					t.Errorf("Start time should be 00:00:00, got %02d:%02d:%02d",
						start.Hour(), start.Minute(), start.Second())
				}

				// Check end is at 23:59:59
				if end.Hour() != 23 || end.Minute() != 59 || end.Second() != 59 {
					t.Errorf("End time should be 23:59:59, got %02d:%02d:%02d",
						end.Hour(), end.Minute(), end.Second())
				}

				// Check same day
				if start.Year() != end.Year() || start.Month() != end.Month() || start.Day() != end.Day() {
					t.Error("Start and end should be on the same day")
				}
			}
		})
	}
}

func TestGetTodayRange(t *testing.T) {
	start, end, err := GetTodayRange("Asia/Shanghai")

	if err != nil {
		t.Fatalf("GetTodayRange failed: %v", err)
	}

	// Check start is before end
	if start.After(end) {
		t.Error("Start time should be before end time")
	}

	// Check duration is about 24 hours
	duration := end.Sub(start)
	expectedDuration := 24 * time.Hour

	if duration > expectedDuration {
		t.Errorf("Duration should be about 24 hours, got %v", duration)
	}
}

func TestParseTimeRange(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		timezone  string
		wantErr   bool
		checkDiff bool
	}{
		{
			name:      "Today keyword",
			input:     "today",
			timezone:  "Asia/Shanghai",
			wantErr:   false,
			checkDiff: true,
		},
		{
			name:      "Yesterday keyword",
			input:     "yesterday",
			timezone:  "Asia/Shanghai",
			wantErr:   false,
			checkDiff: true,
		},
		{
			name:     "Single date",
			input:    "2026-02-11",
			timezone: "Asia/Shanghai",
			wantErr:  false,
		},
		{
			name:     "Date range",
			input:    "2026-02-10,2026-02-12",
			timezone: "Asia/Shanghai",
			wantErr:  false,
		},
		{
			name:     "Invalid date format",
			input:    "invalid-date",
			timezone: "Asia/Shanghai",
			wantErr:  true,
		},
		{
			name:     "Invalid timezone",
			input:    "today",
			timezone: "Invalid/Timezone",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := ParseTimeRange(tt.input, tt.timezone)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimeRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check start is before end
				if start.After(end) {
					t.Error("Start time should be before end time")
				}

				if tt.checkDiff {
					// For today/yesterday, duration should be about 24 hours
					duration := end.Sub(start)
					expectedDuration := 24 * time.Hour
					if duration > expectedDuration {
						t.Errorf("Duration should be about 24 hours, got %v", duration)
					}
				}
			}
		})
	}
}

func TestParseTimeRange_SpecificDate(t *testing.T) {
	start, end, err := ParseTimeRange("2026-02-11", "Asia/Shanghai")

	if err != nil {
		t.Fatalf("ParseTimeRange failed: %v", err)
	}

	// Check both times are on the same day
	if start.Year() != 2026 || start.Month() != 2 || start.Day() != 11 {
		t.Errorf("Start should be on 2026-02-11, got %04d-%02d-%02d",
			start.Year(), start.Month(), start.Day())
	}

	if end.Year() != 2026 || end.Month() != 2 || end.Day() != 11 {
		t.Errorf("End should be on 2026-02-11, got %04d-%02d-%02d",
			end.Year(), end.Month(), end.Day())
	}
}

func TestParseTimeRange_DateRange(t *testing.T) {
	start, end, err := ParseTimeRange("2026-02-10,2026-02-12", "Asia/Shanghai")

	if err != nil {
		t.Fatalf("ParseTimeRange failed: %v", err)
	}

	// Check start is Feb 10
	if start.Year() != 2026 || start.Month() != 2 || start.Day() != 10 {
		t.Errorf("Start should be on 2026-02-10, got %04d-%02d-%02d",
			start.Year(), start.Month(), start.Day())
	}

	// Check end is Feb 12
	if end.Year() != 2026 || end.Month() != 2 || end.Day() != 12 {
		t.Errorf("End should be on 2026-02-12, got %04d-%02d-%02d",
			end.Year(), end.Month(), end.Day())
	}

	// Check start is at 00:00:00
	if start.Hour() != 0 || start.Minute() != 0 || start.Second() != 0 {
		t.Errorf("Start should be at 00:00:00, got %02d:%02d:%02d",
			start.Hour(), start.Minute(), start.Second())
	}

	// Check end is at 23:59:59
	if end.Hour() != 23 || end.Minute() != 59 || end.Second() != 59 {
		t.Errorf("End should be at 23:59:59, got %02d:%02d:%02d",
			end.Hour(), end.Minute(), end.Second())
	}
}
