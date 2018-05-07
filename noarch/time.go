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
type TimeT int32

// NullToTimeT converts a NULL to an array of TimeT.
func NullToTimeT(i int32) []TimeT {
	return []TimeT{}
}

// Time returns the current time.
func Time(tloc []TimeT) TimeT {
	var t = TimeT(int32(time.Now().Unix()))

	if len(tloc) > 0 {
		tloc[0] = t
	}

	return t
}

// IntToTimeT converts an int32 to a TimeT.
func IntToTimeT(t int32) TimeT {
	return TimeT(t)
}

// Ctime converts TimeT to a string.
func Ctime(tloc []TimeT) []byte {
	if len(tloc) > 0 {
		var t = time.Unix(int64(tloc[0]), 0)
		return []byte(t.Format(time.ANSIC) + "\n")
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
	TmSec   int
	TmMin   int
	TmHour  int
	TmMday  int
	TmMon   int
	TmYear  int
	TmWday  int
	TmYday  int
	TmIsdst int
	// tm_gmtoff int32
	// tm_zone   []byte
}

// LocalTime - Convert time_t to tm as local time
// Uses the value pointed by timer to fill a tm structure with the values that
// represent the corresponding time, expressed for the local timezone.
func LocalTime(timer []TimeT) (tm []Tm) {
	t := time.Unix(int64(timer[0]), 0)
	tm = make([]Tm, 1)
	tm[0].TmSec = t.Second()
	tm[0].TmMin = t.Minute()
	tm[0].TmHour = t.Hour()
	tm[0].TmMday = t.Day()
	tm[0].TmMon = int(t.Month()) - 1
	tm[0].TmYear = t.Year() - 1900
	tm[0].TmWday = int(t.Weekday())
	tm[0].TmYday = t.YearDay() - 1
	return
}

// Gmtime - Convert time_t to tm as UTC time
func Gmtime(timer []TimeT) (tm []Tm) {
	t := time.Unix(int64(timer[0]), 0)
	t = t.UTC()
	tm = make([]Tm, 1)
	tm[0].TmSec = t.Second()
	tm[0].TmMin = t.Minute()
	tm[0].TmHour = t.Hour()
	tm[0].TmMday = t.Day()
	tm[0].TmMon = int(t.Month()) - 1
	tm[0].TmYear = t.Year() - 1900
	tm[0].TmWday = int(t.Weekday())
	tm[0].TmYday = t.YearDay() - 1
	return
}

// Mktime - Convert tm structure to time_t
// Returns the value of type time_t that represents the local time described
// by the tm structure pointed by timeptr (which may be modified).
func Mktime(tm []Tm) TimeT {
	t := time.Date(tm[0].TmYear+1900, time.Month(tm[0].TmMon)+1, tm[0].TmMday,
		tm[0].TmHour, tm[0].TmMin, tm[0].TmSec, 0, time.Now().Location())

	tm[0].TmWday = int(t.Weekday())

	return TimeT(int32(t.Unix()))
}

// constants for asctime
var wdayName = [...]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
var monName = [...]string{
	"Jan", "Feb", "Mar", "Apr", "May", "Jun",
	"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
}

// Asctime - Convert tm structure to string
func Asctime(tm []Tm) []byte {
	return []byte(fmt.Sprintf("%.3s %.3s%3d %.2d:%.2d:%.2d %d\n",
		wdayName[tm[0].TmWday],
		monName[tm[0].TmMon],
		tm[0].TmMday, tm[0].TmHour,
		tm[0].TmMin, tm[0].TmSec,
		1900+tm[0].TmYear))
}
