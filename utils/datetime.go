package utils

import "time"

const formatYYYYMMDDHHmmss = "2006-01-02 15:04:05"
const formatYYYYMMDD = "2006-01-02"

func TimeToStringDateTime(t time.Time) string {
	return t.Format(formatYYYYMMDDHHmmss)
}

func TimeToStringDate(t time.Time) string {
	return t.Format(formatYYYYMMDD)
}

func GetDayOfWeek(y, m, d int) time.Weekday {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local).Weekday()
}
