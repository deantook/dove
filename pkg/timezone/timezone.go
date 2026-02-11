package timezone

import (
	"fmt"
	"os"
	"time"
)

// TimezoneConfig 时区配置
type TimezoneConfig struct {
	Timezone string `yaml:"timezone"`
}

// SetTimezone 设置系统时区
func SetTimezone(timezone string) error {
	// 验证时区是否有效
	if _, err := time.LoadLocation(timezone); err != nil {
		return fmt.Errorf("invalid timezone: %s, error: %v", timezone, err)
	}

	// 设置环境变量
	if err := os.Setenv("TZ", timezone); err != nil {
		return fmt.Errorf("failed to set TZ environment variable: %v", err)
	}

	// 设置时区
	time.Local = time.UTC // 先重置为 UTC
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("failed to load timezone location: %v", err)
	}
	time.Local = loc

	return nil
}

// GetCurrentTimezone 获取当前时区
func GetCurrentTimezone() string {
	return time.Local.String()
}

// GetCurrentTime 获取当前时间（使用系统时区）
func GetCurrentTime() time.Time {
	return time.Now()
}

// GetCurrentTimeInTimezone 获取指定时区的当前时间
func GetCurrentTimeInTimezone(timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone: %s, error: %v", timezone, err)
	}
	return time.Now().In(loc), nil
}

// ParseTimeInTimezone 在指定时区解析时间字符串
func ParseTimeInTimezone(timeStr, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone: %s, error: %v", timezone, err)
	}

	// 尝试多种时间格式
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02",
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, timeStr, loc); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time string: %s", timeStr)
}

// FormatTimeInTimezone 格式化时间为指定时区的字符串
func FormatTimeInTimezone(t time.Time, timezone, format string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf("invalid timezone: %s, error: %v", timezone, err)
	}

	return t.In(loc).Format(format), nil
}

// GetAvailableTimezones 获取可用的时区列表
func GetAvailableTimezones() []string {
	return []string{
		"UTC",
		"Asia/Shanghai",
		"Asia/Tokyo",
		"Asia/Seoul",
		"Asia/Singapore",
		"Asia/Hong_Kong",
		"Asia/Bangkok",
		"Asia/Jakarta",
		"Asia/Kolkata",
		"Asia/Dubai",
		"Europe/London",
		"Europe/Paris",
		"Europe/Berlin",
		"Europe/Rome",
		"Europe/Moscow",
		"America/New_York",
		"America/Chicago",
		"America/Denver",
		"America/Los_Angeles",
		"America/Toronto",
		"America/Vancouver",
		"Australia/Sydney",
		"Australia/Melbourne",
		"Pacific/Auckland",
	}
}

// IsValidTimezone 验证时区是否有效
func IsValidTimezone(timezone string) bool {
	_, err := time.LoadLocation(timezone)
	return err == nil
}
