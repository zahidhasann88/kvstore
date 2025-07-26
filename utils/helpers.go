package utils

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return d.String()
	}

	seconds := int(d.Seconds())
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / 60
	remainingSeconds := seconds % 60

	if minutes < 60 {
		if remainingSeconds > 0 {
			return fmt.Sprintf("%dm%ds", minutes, remainingSeconds)
		}
		return fmt.Sprintf("%dm", minutes)
	}

	hours := minutes / 60
	remainingMinutes := minutes % 60

	if remainingMinutes > 0 {
		return fmt.Sprintf("%dh%dm", hours, remainingMinutes)
	}
	return fmt.Sprintf("%dh", hours)
}

func IsValidKey(key string) bool {
	return len(key) > 0 && len(key) <= 250
}
