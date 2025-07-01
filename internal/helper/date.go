package helper

import (
	"time"
)

func DateNowJakarta() string {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	return time.Now().In(loc).Format("2006-01-02")
}
func DatetimeNowJakartaString() string {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}

func AddMinutes(m int) string {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return time.Now().In(loc).Add(time.Minute * time.Duration(m)).Format("2006-01-02 15:04:05")
}

func DateFilename() string {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return time.Now().In(loc).Format("20060102150405")
}

func DatetimeNowJakarta() time.Time {
	jakartaLocation, _ := time.LoadLocation("Asia/Jakarta")

	return time.Now().In(jakartaLocation)
}

func ParseDate(dateStr string) (time.Time, error) {
	const layout = "2006-01-02"
	return time.Parse(layout, dateStr)
}
