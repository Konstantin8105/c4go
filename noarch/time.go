package noarch

import (
	"fmt"
	"time"
)

// TimeT is the representation of "time_t".
// For historical reasons, it is generally implemented as an integral value
// representing the number of seconds elapsed
// since 00:00 hours, Jan 1, 1970 UTC (i.e., a unix timestamp).
// Although libraries may implement this type using alternative time
// representations.
type TimeT int

// NullToTimeT converts a NULL to an array of TimeT.
func NullToTimeT(i int) *TimeT {
	return nil
}

// Time returns the current time.
func Time(tloc *TimeT) TimeT {
	var t = TimeT(int(time.Now().Unix()))

	if tloc != nil {
		*tloc = t
	}

	return t
}

// IntToTimeT converts an int to a TimeT.
func IntToTimeT(t int) TimeT {
	return TimeT(t)
}

// Ctime converts TimeT to a string.
func Ctime(tloc *TimeT) *byte {
	if tloc != nil {
		var t = time.Unix(int64(*tloc), 0)
		return &append([]byte(t.Format(time.ANSIC)+"\n"), 0)[0]
	}

	return nil
}

// TimeTToFloat64 converts TimeT to a float64. It is used by the tests.
func TimeTToFloat64(t TimeT) float64 {
	return float64(t)
}

// Tm - base struct in "time.h"
// Structure containing a calendar date and time broken down into its
// components
type Tm struct {
	Tm_sec   int
	Tm_min   int
	Tm_hour  int
	Tm_mday  int
	Tm_mon   int
	Tm_year  int
	Tm_wday  int
	Tm_yday  int
	Tm_isdst int
	// tm_gmtoff int
	// tm_zone   []byte
}

// Localtime - Convert time_t to tm as local time
// Uses the value pointed by timer to fill a tm structure with the values that
// represent the corresponding time, expressed for the local timezone.
func LocalTime(timer *TimeT) (tm *Tm) {
	t := time.Unix(int64(*timer), 0)
	tm = &Tm{}
	tm.Tm_sec = int(t.Second())
	tm.Tm_min = int(t.Minute())
	tm.Tm_hour = int(t.Hour())
	tm.Tm_mday = int(t.Day())
	tm.Tm_mon = int(t.Month() - 1)
	tm.Tm_year = int(t.Year() - 1900)
	tm.Tm_wday = int(t.Weekday())
	tm.Tm_yday = int(t.YearDay() - 1)
	return
}

// Gmtime - Convert time_t to tm as UTC time
func Gmtime(timer *TimeT) (tm *Tm) {
	t := time.Unix(int64(*timer), 0)
	t = t.UTC()
	tm = &Tm{}
	tm.Tm_sec = int(t.Second())
	tm.Tm_min = int(t.Minute())
	tm.Tm_hour = int(t.Hour())
	tm.Tm_mday = int(t.Day())
	tm.Tm_mon = int(t.Month() - 1)
	tm.Tm_year = int(t.Year() - 1900)
	tm.Tm_wday = int(t.Weekday())
	tm.Tm_yday = int(t.YearDay() - 1)
	return
}

// Mktime - Convert tm structure to time_t
// Returns the value of type time_t that represents the local time described
// by the tm structure pointed by timeptr (which may be modified).
func Mktime(tm *Tm) TimeT {
	t := time.Date(int(tm.Tm_year+1900), time.Month(tm.Tm_mon)+1, int(tm.Tm_mday),
		int(tm.Tm_hour), int(tm.Tm_min), int(tm.Tm_sec), 0, time.Now().Location())

	tm.Tm_wday = int(t.Weekday())

	return TimeT(int(t.Unix()))
}

// constants for asctime
var wday_name = [...]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
var mon_name = [...]string{
	"Jan", "Feb", "Mar", "Apr", "May", "Jun",
	"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
}

// Asctime - Convert tm structure to string
func Asctime(tm *Tm) *byte {
	return &append([]byte(fmt.Sprintf("%.3s %.3s%3d %.2d:%.2d:%.2d %d\n",
		wday_name[tm.Tm_wday],
		mon_name[tm.Tm_mon],
		tm.Tm_mday, tm.Tm_hour,
		tm.Tm_min, tm.Tm_sec,
		1900+tm.Tm_year)), 0)[0]
}
