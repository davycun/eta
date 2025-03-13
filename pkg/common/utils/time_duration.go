package utils

import (
	"fmt"
	"time"
)

func FmtDuration(d time.Duration, hour, minute string) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	if h <= 0 {
		return fmt.Sprintf("%d%s", m, minute)
	}
	return fmt.Sprintf("%d%s%d%s", h, hour, m, minute)
}
