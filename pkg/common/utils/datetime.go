package utils

import (
	"errors"
	"fmt"
	"time"
)

var (
	timeFormat = []string{
		"2006-01-02T15:04:05+08:00",
		"2006-01-02T15:04:05.000000+08:00",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02 15",
		"2006-01-02",
		"2006-01",
		"01-02",
		"01-02-06",
		"02-01-06 15:04:05",
		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02 15",
		"2006/01/02",
		"2006/1/2",
		"2006/01",
		"01/02",
		"02/01/06 15:04:05",
		"20060102",
		"010206",
		"2006",
		"06",
		"01",
		"15:04:05",
		"15:04",
		"04:05",
	}
)

// FormatStrToTime 字符串转时间。时区为当前系统时区
func FormatStrToTime(str string, format ...string) (time.Time, error) {
	if len(format) == 0 {
		format = timeFormat
	}
	for _, f := range format {
		toTime, err := FormatStrToTimeWithTz(str, f, "Local")
		if err == nil {
			return toTime, nil
		}
	}
	return time.Time{}, errors.New(fmt.Sprintf("解析异常: %s", str))
}

func FormatStrToTimeWithTz(str, format string, timezone ...string) (time.Time, error) {
	if timezone != nil && timezone[0] != "" {
		loc, err := time.LoadLocation(timezone[0])
		if err != nil {
			return time.Time{}, err
		}

		return time.ParseInLocation(format, str, loc)
	}
	return time.Parse(format, str)
}
