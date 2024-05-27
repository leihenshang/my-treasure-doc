package utils_datetime

import "time"

func TimeToFormat(time *time.Time) string {
	if time == nil {
		return ""
	}

	return time.Format("2006-01-02 15:04:05")
}
