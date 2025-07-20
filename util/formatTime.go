package util

import "time"

func FormatTime(t time.Time) string {
	layout := "2006-01-02 15:04:05"
	return t.Format(layout)
}
