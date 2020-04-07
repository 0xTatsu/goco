package pkg

import (
	"time"
)

const (
	DateTimeFormatMySQL = "2006-01-02 15:04:05"
)

type Time struct {
	Location *time.Location
}

func NewTime(location *time.Location) *Time {
	return &Time{
		Location: location,
	}
}

func (r *Time) CurrentTimeInMySQLFormat() string {
	return time.Now().In(r.Location).Format(DateTimeFormatMySQL)
}
