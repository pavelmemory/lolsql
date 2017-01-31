package types

import (
	"time"
	"fmt"
	"database/sql/driver"
)

type CustomNullableUserType interface {
	CustomUserType
	IsNull() bool
}

type CustomUserType interface {
	// This method used to initialise Go object with value retrieved from database
	// Receive must be pointer.
	// func (this *T) Scan(src interface{}) (err error) {...}
	Scan(src interface{}) (err error)
	// This method used to convert Go object to database representation.
	// Supported types to return are: int64, float64, bool, []byte, string, time.Time
	// Receive must be value.
	// func (this T) Value() (driver.Value, error) {...}
	Value() (driver.Value, error)
}

const timeFormat = "2006-01-02 15:04:05.999999"

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (this *NullTime) Scan(src interface{}) (err error) {
	switch v := src.(type) {
	case time.Time:
		this.Time, this.Valid = v, true
	case []byte:
		this.Time, err = parseDateTime(string(v), time.UTC)
		this.Valid = (err == nil)
	case string:
		this.Time, err = parseDateTime(v, time.UTC)
		this.Valid = (err == nil)
	case nil:
		this.Valid = false
	default:
		this.Valid = false
		err = fmt.Errorf("Can't convert %T to time.Time", src)
	}
	return
}

var parselength = map[int]bool {
	10:true, 19:true, 21:true, 22:true, 23:true, 24:true, 25:true, 26:true,  // up to "YYYY-MM-DD HH:MM:SS.MMMMMM"
}

func parseDateTime(str string, loc *time.Location) (time.Time, error) {
	t := time.Time{}
	if parselength[len(str)] {
		if str == "0000-00-00 00:00:00.0000000"[:len(str)] {
			return t, nil
		}
		t, err := time.Parse(timeFormat[:len(str)], str)
		// Adjust location
		if err == nil && loc != time.UTC {
			y, mo, d := t.Date()
			h, mi, s := t.Clock()
			t, err = time.Date(y, mo, d, h, mi, s, t.Nanosecond(), loc), nil
		}
		return t, err
	}
	return t, fmt.Errorf("Invalid time string: %s", str)
}
