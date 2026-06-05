package helper

import (
	"fmt"
	"time"
)

func DateNowJakarta() string {
	return time.Now().In(time.Local).Format("2006-01-02")
}

func DatetimeNowJakartaString() string {
	return time.Now().In(time.Local).Format("2006-01-02 15:04:05")
}

func AddMinutes(m int) string {
	return time.Now().In(time.Local).Add(time.Minute * time.Duration(m)).Format("2006-01-02 15:04:05")
}

func DateFilename() string {
	return time.Now().In(time.Local).Format("20060102150405")
}

func DatetimeNowJakarta() time.Time {
	return time.Now().In(time.Local)
}

func ParseDate(dateStr string) (time.Time, error) {
	const layout = "2006-01-02"
	return time.Parse(layout, dateStr)
}

func NowStrUTC() string {
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(),
		time.Now().UTC().Hour(), time.Now().UTC().Minute(), time.Now().UTC().Second())
}

func ConvertToJakartaTime(t time.Time) string {
	return t.In(time.Local).Format("2006-01-02 15:04:05")
}

func ConvertToJakartaDate(t time.Time) string {
	return t.In(time.Local).Format("2006-01-02")
}
