package otime

import (
	"fmt"
	"strings"
	"time"
)

type OTime struct {
	time.Time
}

const (
	layoutTime            = "15:04"
	layoutTimeWithSeconds = "15:04:05"
	layoutDate            = "20060102"
)

var nilTime = (time.Time{}).UnixNano()

// DefaultLocation is the default timezone to be used when creating a new standup config.
// TODO will need to remove this and use channel-specific location
// when adding per-channel timezone setting
var DefaultLocation *time.Location

func Parse(value string) (OTime, error) {
	argTime, err := time.Parse(layoutTime, value)
	if err != nil {
		return OTime{}, err
	}

	now := time.Now()
	argTime = time.Date(now.Year(), now.Month(), now.Day(), argTime.Hour(), argTime.Minute(), 0, 0, DefaultLocation)
	return OTime{argTime}, nil
}

func Now(timezone string) OTime {
	now := time.Now()
	location, _ := time.LoadLocation(timezone)
	return OTime{now.In(location)}
}

// GetTime returns time with format like "15:04"
func (ct OTime) GetTime(timezone string) OTime {
	now, _ := time.Parse(layoutTime, ct.Format(layoutTime))
	location, _ := time.LoadLocation(timezone)
	return OTime{now.In(location)}
}

// GetTimeWithSeconds returns time with format like "15:04:05"
func (ct OTime) GetTimeWithSeconds(timezone string) OTime {
	now, _ := time.Parse(layoutTimeWithSeconds, ct.Format(layoutTimeWithSeconds))
	location, _ := time.LoadLocation(timezone)
	return OTime{now.In(location)}
}

func (ct OTime) GetTimeString() string {
	return ct.Time.Format(layoutTime)
}

// GetDate returns date with format like "20060102"
func (ct OTime) GetDate(timezone string) OTime {
	now, _ := time.Parse(layoutDate, ct.Format(layoutDate))
	location, _ := time.LoadLocation(timezone)
	return OTime{now.In(location)}
}

func (ct OTime) GetDateString() string {
	return ct.Time.Format(layoutDate)
}

func (ct *OTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	t, err := Parse(s)
	if err != nil {
		return err
	}

	ct.Time = t.Time
	return
}

func (ct OTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(layoutTime))), nil
}
