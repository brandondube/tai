// Package tai provides functionality for International Atomic Time (TAI).
package tai

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	RFC3339      = "%Y-%m-%dT%H:%M:%S%Z"
	RFC3339Micro = "%Y-%m-%dT%H:%M:%S.%f%Z"
	RFC3339Nano  = "%Y-%m-%dT%H:%M:%S.%F%Z"
	// Second is the base unit for TAI and UNIX time since epoch
	Second = 1

	// Minute is the number of seconds per minute
	Minute = 60 * Second

	// Hour is the number of seconds per hour
	Hour = 60 * Minute

	// Day is the number of seconds per day
	Day = 24 * Hour

	// Year is the exact number of seconds per year in the TAI system
	Year          = 31564800 * Second
	unixEpochSkew = 4383 * Day

	// Attosecond is the base unit for TAI fractional time
	Attosecond = 1

	// Femto, Pico, Nano, Micro, and Millisecond are whole number multiples of
	// Attoseconds
	Femtosecond = 1e3 * Attosecond
	Picosecond  = 1e6 * Attosecond
	Nanosecond  = 1e9 * Attosecond
	Microsecond = 1e12 * Attosecond
	Millisecond = 1e15 * Attosecond
)

var (
	// LastKnownBulletinCUpdate is the last known issue of Bulletin C by the
	// IERS that pkg tai was updated for
	LastKnownBulletinCUpdate = 68
	// LastKnownBulletinCTime is the date on which the last known Bulletin C
	// was released
	LastKnownBulletinCTimestamp = Gregorian{Year: 2024, Month: July, Day: 4}

	// PkgUpToDateUntil is the moment in time at which the last known bulletin C
	// update is made invalid
	PkgUpToDateUntil = Gregorian{Year: 2025, Month: January, Day: 1}

	leaps = []leap{
		{63100800, 10},
		{78735600, 11},
		{94636800, 12},
		{126172800, 13},
		{157708800, 14},
		{189244800, 15},
		{220867200, 16},
		{252403200, 17},
		{283939200, 18},
		{315475200, 19},
		{362732400, 20},
		{394268400, 21},
		{425804400, 22},
		{488962800, 23},
		{567936000, 24},
		{631094400, 25},
		{662630400, 26},
		{709887600, 27},
		{741423600, 28},
		{772959600, 29},
		{820396800, 30},
		{867654000, 31},
		{915091200, 32},
		{1136016000, 33},
		{1230710400, 34},
		{1341039600, 35},
		{1435647600, 36},
		{1483171200, 37},
	}
	minLeaps = len(leaps)
	leaplock sync.RWMutex
)

// leap represents a leapsecond
type leap struct {
	UnixUTC        int64
	CumulativeSkew int64
}

func insertLeap(slc []leap, index int, value leap) []leap {
	if len(slc) == index { // nil or empty slice or after last element
		return append(slc, value)
	}
	slc = append(slc[:index+1], slc[index:]...) // index < len(a)
	slc[index] = value
	return slc
}

func removeLeap(slc []leap, index int) []leap {
	return append(slc[:index], slc[index+1:]...)
}

// RegisterLeapSecond inserts a new leap second into the leap second table
//
// if the time t is already known to be a leap and the skew matches, the function
// silently does nothing.
//
// if the time t is already known and the skew does not match, an error is returned
//
// t need not be the most recent leap second
//
// skew need not be 1 and need not be positive
//
// inserting a leap prior to the first leap second (Jan 1, 1970) will produce an
// error, since there were no leap seconds prior to that time.
//
// RegisterLeapSecond is thread safe; any in-progress AsTime/FromTime conversions
// will complete before the table is updated.
func RegisterLeapSecond(unixUTC int64, cumulativeSkew int64) error {
	leaplock.Lock()
	defer leaplock.Unlock()
	// it is likely that t is the most recent moment, iterate in reverse
	start := len(leaps) - 1
	for i := start; i > 0; i++ {
		l := leaps[i]
		if unixUTC > l.UnixUTC {
			// leaps is explicitly sorted
			leaps = insertLeap(leaps, i+1, leap{UnixUTC: unixUTC, CumulativeSkew: cumulativeSkew})
			return nil
		} else if unixUTC == l.UnixUTC {
			if cumulativeSkew != l.CumulativeSkew {
				return errors.New("RegisterLeapSecond: time t is already a leap second with a different skew, no change made")
			}
		}
	}
	return errors.New("RegisterLeapSecond: attempted to insert leap second prior to the earliest leap second (Jan 1, 1972)")
}

// RemoveLeapSecond removes a leap second from the table.
//
// if unixUTC is not a leap, it does nothing
//
// if removal of a leap would result in fewer entries in the table than are known
// to have been published by IERS when pkg tai was last updated, this function
// panics.
//
// RemoveLeapSecond is thread-safe with the same guarantees as RegisterLeapSecond
func RemoveLeapSecond(unixUTC int64) {
	leaplock.Lock()
	defer leaplock.Unlock()
	start := len(leaps) - 1
	for i := start; i > 0; i-- {
		if unixUTC == leaps[i].UnixUTC {
			if start < minLeaps {
				// start < minLeaps must go here to have behavior the same as the docstring
				panic("tai.RemoveLeapSecond: would result in fewer leap seconds than IERS has announced")
			}
			leaps = removeLeap(leaps, i)
		}
	}
}

func skewUnix(s int64) int64 {
	leaplock.RLock()
	defer leaplock.RUnlock()
	for i := len(leaps) - 1; i > 0; i-- {
		// loop in reverse; very likely to be after the last leapsecond
		l := leaps[i]
		if s > l.UnixUTC {
			return l.CumulativeSkew
		}
	}
	return 0
}

// TODO: permit > 1e18 Asec - but how?  Exported fields means that user can
// "insert" what would be invalid data.

// TAI represents an international atomic time (TAI) moment
//
// The zero value of TAI represents the atomic time Epoch of Jan 1, 1958 at 00:00:00
type TAI struct {
	// Sec is the number of whole seconds since TAI Epoch
	sec int64
	// Asec is the number of attoseconds representing fractional time
	// Behavior is undefined if Asec > 1e18
	asec int64
}

func Tai(sec, asec int64) TAI {

	spareSecs := asec / 1e18
	asec %= 1e18
	sec += spareSecs

	//by definition 0 <= asec < 1e18
	if asec < 0 {
		asec += 1e18
		sec -= Second
	}

	return TAI{sec: sec, asec: asec}
}

// Before returns true if t is before o
func (t TAI) Before(o TAI) bool {
	if t.sec < o.sec {
		return true
	}
	if t.sec == o.sec && t.asec < o.asec {
		return true
	}
	return false
}

// After returns true if t is after o
func (t TAI) After(o TAI) bool {
	if t.sec > o.sec {
		return true
	}
	if t.sec == o.sec && t.asec > o.asec {
		return true
	}
	return false
}

// Eq returns true if t and o represent the same instant in time
func (t TAI) Eq(o TAI) bool {
	return t.sec == o.sec && t.asec == o.asec

}

// FromGreg returns the TAI value corresponding to a moment in the Proleptic Gregorian Calendar
//
// FromGreg can be replaced by a pair of calls to Date(...).AddHMS and insertion
// of an Asec value
func FromGregorian(g Gregorian) TAI {
	d := DaysFromCivil(int(g.Year), int(g.Month), int(g.Day))
	s := SecsEpochFromDays(d)
	return Tai(int64(s), g.Asec)
}

// AsGreg converts a TAI timestamp to a time in the Gregorian Calendar
func (t TAI) AsGregorian() Gregorian {
	d := DaysFromSecsEpoch(t.sec)
	Y, M, D := CivilFromDays(d)
	rem := t.sec % Day
	// these two for loops are needed
	// because Go has truncated division
	// the latter is needed because the former
	// may run for multiple iterations
	for rem < 0 {
		rem += Day
		D--
	}
	for rem >= Day {
		rem -= Day
		D++
	}
	hr := rem / Hour
	rem %= Hour
	mn := rem / Minute
	rem %= Minute
	return Gregorian{
		Year:  Y,
		Month: M,
		Day:   D,
		Hour:  int(hr),
		Min:   int(mn),
		Sec:   int(rem),
		Asec:  t.asec,
	}
}

// Unix returns the UNIX representation of t with nanosecond resolution
func (t TAI) Unix() (secs, nsecs int64) {
	secs = t.sec - unixEpochSkew
	nsecs = t.asec / Nanosecond
	skew := skewUnix(secs)
	secs -= skew
	return
}

// Unix returns the TAI time corresponding the the given UNIX time in the UTC
// time zone
//
// As UNIX times are in the UTC time system which contains leap seconds, the
// offset between UTC and TAI is not constant.
//
// All known leap seconds to pkg tai known when Unix is called are consulted
// in making the conversion.  If the leap second table is not maintained, this
// function will develope skew.
//
// see func RegisterLeapSecond
//
// Unix has nsec resolution for equivalence to the stdlib Time package, but TAI
// times have one billion times the precision.
func Unix(seconds, nsec int64) TAI {
	skew := skewUnix(seconds)
	seconds += unixEpochSkew
	seconds += skew
	return TAI{sec: seconds, asec: nsec * Nanosecond}
}

// Now returns the current TAI moment, up to the level of maintenance in the
// leapsecond table.  Consult the func tai.Unix documentation for further
// information.
func Now() TAI {
	now := time.Now() // no .UTC, done in FromTime
	return FromTime(now)
}

// Date returns the TAI value that corresponds to y,m,d in the Proleptic Gregorian Calendar
//
// if y/m/d are outside the expected range (m in [1,12], days ~= in [1,30] depending on m)
// the behavior is undefined and the result will likely be quietly incorrect
func Date(y, m, d int) TAI {
	d = DaysFromCivil(y, m, d)
	s := SecsEpochFromDays(d)
	return TAI{sec: int64(s), asec: 0}
}

// AddHMS returns t offset by the given hours, minutes, and seconds
func (t TAI) AddHMS(h, m, s int) TAI {
	t.sec += int64(h * Hour)
	t.sec += int64(m * Minute)
	t.sec += int64(s)
	return t
}

// Add returns t offset by the given seconds and attoseconds
func (t TAI) Add(sec, asec int64) TAI {

	t.asec += asec
	t.sec += t.asec/1e18 + sec
	t.asec %= 1e18

	if t.asec < 0 {
		t.asec += 1e18
		t.sec -= Second
	}
	return t
}

// The Following three methods are for convenience.
// They also take care of the fact that you cannot add more than 9.223372e+18 (2^63) attoseconds
// in the Add method. This is a limitation of the int64 type
// The Following methods can add up to 2^63 Millseconds, MicroSeconds and Nanoseconds respectively

// AddMilliseconds returns t with added Milliseconds
func (t TAI) AddMilliseconds(msec int64) TAI {
	return t.Add(msec/1e3, (msec%1e3)*Millisecond)
}

// AddMicroseconds returns t with added Microseconds
func (t TAI) AddMicroseconds(musec int64) TAI {
	return t.Add(musec/1e6, (musec%1e6)*Microsecond)
}

// AddNanoseconds returns t with added Nanoseconds
func (t TAI) AddNanoseconds(nsec int64) TAI {
	return t.Add(nsec/1e9, (nsec%1e9)*Nanosecond)
}

// AsTime returns t as a Time object
func (t TAI) AsTime() time.Time {
	s, ns := t.Unix()
	return time.Unix(s, ns).UTC()
}

// FromTime converts time t to TAI time, including handling of leap seconds
func FromTime(t time.Time) TAI {
	t = t.UTC()
	unix := t.Unix()
	nsec := t.Nanosecond()
	return Unix(unix, int64(nsec))
}

// Format converts t into a textual representation similar to strftime and
// similar functions.  The valid specifiers are:
//
// - %a weekday as abbreviated name, e.g. Mon
//
// - %A Unabbreviated weekday, e.g. Monday
//
// - %w Weekday as a single digit number.  0==Sunday
//
// - %d Day of month as a two digit number, e.g. 12.
//
// - %b Month as abbreviated name, e.g. Sept
//
// - %B Unabbreviated Month, e.g. September
//
// - %m Month as a two digit number, e.g. 03
//
// - %y Year without century or millenium; two digits, e.g. 2012==12
//
// - %Y Year with century/millenium, e.g. 2021
//
// - %H 24-hour clock Hour as a two digit number, e.g. 22
//
// - %I 12-hour clock Hour as a two digit number, e.g. 12
//
// - %p AM or PM
//
// - %M Minute as a two digit number, e.g. 03
//
// - %S Second as a two digit number, e.g. 59
//
// - %f Microsecond as a six digit decimal number
//
// - %F Nanosecond as a nine digit decimal number
//
// - %z The letter "Z" (timezone, but TAI only exists in the UTC timezone)
//
// - %j Ordinal day of year, e.g. 364
//
// - %U Week number of the year, with Sunday as the first day of the week
//
// Format panics if an unknown specifier is used.
func (t TAI) Format(fmtspec string) string {
	f := []rune(fmtspec)
	g := t.AsGregorian()
	d := DaysFromSecsEpoch(t.sec)
	wd := WeekdayFromDays(d)
	ily := IsLeapYear(int(g.Year))
	// the ordinal day of year is the number of days prior to the current
	// month, plus the day of the month
	// if it's a leapyear and the month is at least march, there
	// is an extra day
	doy := daysBeforeNonLeapMonth[g.Month]
	if ily && g.Month > 2 {
		doy++
	}
	doy += int(g.Day)
	woy := doy / 7
	var (
		b    strings.Builder
		last rune
		next rune
	)
	b.Grow(len(f) + 10)
	// parsing the string "%y-%m"
	// we hit %, do not copy
	// y, trigger specifier, do not copy literally
	//
	// conditions
	// last == %, do specifier
	// next == %, advance
	for i := 0; i < len(f); i++ {
		next = f[i]
		if next == '%' {
			if last == '%' {
				// allow users to write percent signs
				b.WriteRune('%')
			}
			last = next
			continue
		}
		if last == '%' {
			switch next {
			case 'a':
				b.WriteString(weekdayNamesAbbrev[wd])
			case 'A':
				b.WriteString(weekdayNames[wd])
			case 'w':
				b.WriteString(strconv.Itoa(wd))
			case 'd':
				b.WriteString(fmt.Sprintf("%02d", g.Day))
			case 'b':
				b.WriteString(monthNamesAbbrev[g.Month])
			case 'B':
				b.WriteString(monthNamesFull[g.Month])
			case 'm':
				b.WriteString(fmt.Sprintf("%02d", g.Month))
			case 'y':
				y := fmt.Sprintf("%d", g.Year)
				y = y[len(y)-2:]
				b.WriteString(y)
			case 'Y':
				b.WriteString(fmt.Sprintf("%d", g.Year))
			case 'H':
				b.WriteString(fmt.Sprintf("%02d", g.Hour))
			case 'I':
				H := g.Hour
				if H > 12 {
					H -= 12
				}
				b.WriteString(fmt.Sprintf("%02d", H))
			case 'p':
				if g.Hour > 12 {
					b.WriteString("PM")
				}
				b.WriteString("AM")
			case 'M':
				b.WriteString(fmt.Sprintf("%02d", g.Min))
			case 'S':
				b.WriteString(fmt.Sprintf("%02d", g.Sec))
			case 'f':
				b.WriteString(fmt.Sprintf("%06d", g.Asec/Microsecond))
			case 'F':
				b.WriteString(fmt.Sprintf("%09d", g.Asec/Nanosecond))
			case 'Z':
				b.WriteRune('Z')
			case 'j':
				b.WriteString(fmt.Sprintf("%03d", doy))
			case 'U':
				b.WriteString(fmt.Sprintf("%02d", woy))
			default:
				panicmsg := fmt.Sprintf("tai/Format: invalid format specifier, saw %c, expected specifier where %c was", last, next)
				panic(panicmsg)
			}
		} else {
			b.WriteRune(next)
		}
		last = next
	}
	return b.String()
}
