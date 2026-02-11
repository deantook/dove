package timezone

import (
	"testing"
	"time"
)

func TestSetTimezone(t *testing.T) {
	tests := []struct {
		name      string
		timezone  string
		expectErr bool
	}{
		{"valid timezone", "Asia/Shanghai", false},
		{"valid timezone UTC", "UTC", false},
		{"invalid timezone", "Invalid/Timezone", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetTimezone(tt.timezone)
			if (err != nil) != tt.expectErr {
				t.Errorf("SetTimezone() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestGetCurrentTimezone(t *testing.T) {
	// 设置一个已知的时区
	err := SetTimezone("Asia/Shanghai")
	if err != nil {
		t.Fatalf("Failed to set timezone: %v", err)
	}

	timezone := GetCurrentTimezone()
	if timezone == "" {
		t.Error("GetCurrentTimezone() returned empty string")
	}
}

func TestGetCurrentTime(t *testing.T) {
	time := GetCurrentTime()
	if time.IsZero() {
		t.Error("GetCurrentTime() returned zero time")
	}
}

func TestGetCurrentTimeInTimezone(t *testing.T) {
	tests := []struct {
		name      string
		timezone  string
		expectErr bool
	}{
		{"valid timezone", "Asia/Shanghai", false},
		{"invalid timezone", "Invalid/Timezone", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time, err := GetCurrentTimeInTimezone(tt.timezone)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetCurrentTimeInTimezone() error = %v, expectErr %v", err, tt.expectErr)
			}
			if !tt.expectErr && time.IsZero() {
				t.Error("GetCurrentTimeInTimezone() returned zero time")
			}
		})
	}
}

func TestParseTimeInTimezone(t *testing.T) {
	tests := []struct {
		name      string
		timeStr   string
		timezone  string
		expectErr bool
	}{
		{"valid datetime", "2023-01-01 12:00:00", "Asia/Shanghai", false},
		{"valid date", "2023-01-01", "Asia/Shanghai", false},
		{"invalid timezone", "2023-01-01 12:00:00", "Invalid/Timezone", true},
		{"invalid time string", "invalid time", "Asia/Shanghai", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTimeInTimezone(tt.timeStr, tt.timezone)
			if (err != nil) != tt.expectErr {
				t.Errorf("ParseTimeInTimezone() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestFormatTimeInTimezone(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	format := "2006-01-02 15:04:05"

	tests := []struct {
		name      string
		timezone  string
		expectErr bool
	}{
		{"valid timezone", "Asia/Shanghai", false},
		{"invalid timezone", "Invalid/Timezone", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FormatTimeInTimezone(testTime, tt.timezone, format)
			if (err != nil) != tt.expectErr {
				t.Errorf("FormatTimeInTimezone() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestGetAvailableTimezones(t *testing.T) {
	timezones := GetAvailableTimezones()
	if len(timezones) == 0 {
		t.Error("GetAvailableTimezones() returned empty slice")
	}

	// 检查是否包含一些常见的时区
	expectedTimezones := []string{"UTC", "Asia/Shanghai", "America/New_York"}
	for _, expected := range expectedTimezones {
		found := false
		for _, tz := range timezones {
			if tz == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected timezone %s not found in available timezones", expected)
		}
	}
}

func TestIsValidTimezone(t *testing.T) {
	tests := []struct {
		name     string
		timezone string
		expected bool
	}{
		{"valid timezone", "Asia/Shanghai", true},
		{"valid timezone UTC", "UTC", true},
		{"invalid timezone", "Invalid/Timezone", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidTimezone(tt.timezone)
			if result != tt.expected {
				t.Errorf("IsValidTimezone() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
