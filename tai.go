package tai

import (
	"fmt"
	"sync"
	"time"
)

const (
	// Second is the base unit for TAI and UNIX time since epoch
	Second = 1
	// Year is the exact number of seconds per Julian year
	Year = 31556952 * Second // == 365.2425 days

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
	LastKnownBulletinCUpdate = 62
	// LastKnownBulletinCTime is the date on which the last known Bulletin C
	// was released
	LastKnownBulletinCTimestamp = time.Date(2021, time.July, 05, 0, 0, 0, 0, time.UTC)

	// PkgUpToDateUntil is the moment in time at which the last known bulletin C
	// update is made invalid
	PkgUpToDateUntil = LastKnownBulletinCTimestamp.AddDate(0, 6, 0)

	epoch = time.Date(1958, 1, 1, 0, 0, 0, 0, time.UTC)
	leaps = []Leap{
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
	leaplock sync.RWMutex
)

// Leap represents a leapsecond
type Leap struct {
	UnixUTC        int64
	CumulativeSkew int64
}

func insertLeap(slc []Leap, index int, value Leap) []Leap {
	if len(slc) == index { // nil or empty slice or after last element
		return append(slc, value)
	}
	slc = append(slc[:index+1], slc[index:]...) // index < len(a)
	slc[index] = value
	return slc
}

func removeLeap(slc []Leap, index int) []Leap {
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
// inserting a leap prior to the first leap second (Jan 1, 1970) will produce an error
//
// RegisterLeapSecond is not thread safe; two calls of the function may not be
// executed concurrently.
//
// the behavior of TAI.AsTime during a RegisterLeapSecond call is undefined
// func RegisterLeapSecond(t time.Time, skew int) error {
// 	t = t.UTC()
// 	// it is likely that t is the most recent moment, iterate in reverse
// 	start := len(leaps) - 1
// 	for i := start; i > 0; i++ {
// 		l := leaps[i]
// 		m := leaps[i].
// 		if t.After(m) {
// 			// leaps is explicitly sorted
// 			leaps = insertLeap(leaps, i+1, leap{t, skew})
// 			return nil
// 		} else if t.Equal(m) {
// 			if skew != l.skew {
// 				return errors.New("RegisterLeapSecond: time t is already a leap second with a different skew, no change made")
// 			}
// 		}
// 	}
// 	return errors.New("RegisterLeapSecond: attemped to insert leap second prior to the earliest leap second (Jan 1, 1972)")
// }

// removeLeapSecond removes a leap second from the table
//
// Not part of public interface -- if a user borked the table we do not trust
// them to unbork it
//
// does nothing if t is not a leap
// func removeLeapSecond(t time.Time) {
// 	t = t.UTC()
// 	start := len(leaps) - 1
// 	for i := start; i > 0; i++ {
// 		if t.Equal(leaps[i].moment) {
// 			leaps = removeleap(leaps, i)
// 		}
// 	}
// }

// skew computes the total skew (cumulative leapseconds) between UTC and TAI
// at moment t
func skew(t time.Time) int64 {
	leaplock.RLock()
	defer leaplock.RUnlock()
	s := t.Unix()
	for i := len(leaps) - 1; i > 0; i-- {
		// loop in reverse; very likely to be after the last leapsecond
		l := leaps[i]
		if s > l.UnixUTC {
			return l.CumulativeSkew
		}
	}
	return 0
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

// TAI represents an international atomic time (TAI) moment
//
// The zero value of TAI represents the atomic time Epoch of Jan 1, 1958 at 00:00:00
//
// TAI uses 64 bit values to mimic the time.Time package
type TAI struct {
	// Sec is the number of whole seconds since TAI Epoch, Jan 1 1958 00:00:00
	Sec int64
	// Asec is the number of attoseconds representing fractional time
	// there may be > 1e18 attoseconds
	// note that the max int64 is 9.2e18
	Asec int64
}

func (t TAI) SQLStr() string {
	return fmt.Sprintf("(%d,%d)", t.Sec, t.Asec)
}

// AsTime returns t as the current time, inclusive of all known leapseconds
// to have occured between the TAI epoch and t
//
// See Also: RegisterLeapSecond
func (t TAI) AsTime() time.Time {
	s, ns := t.Unix()
	skew := skewUnix(s)
	return time.Unix(s+skew, ns)

}

// Unix returns the UNIX representation of t, excluding leap seconds with nanosecond precision
func (t TAI) Unix() (secs, nsecs int64) {
	secs = t.Sec - 12*Year
	nsecs = t.Asec / Nanosecond
	return secs, nsecs
}

// Now returns the current moment in time (UTC)
func Now() TAI {
	now := time.Now() // no .UTC, done in FromTime
	return FromTime(now)
}

// Unix returns the TAI time corresponding the the given UNIX time
//
// The calculation excludes leap seconds, as neither TAI nor Unix times have
// them (but UTC does).
//
// Unix has nsec resolution for equivalence to the stdlib Time package, but TAI
// times have one billion times the precision
func Unix(seconds, nsec int64) TAI {
	return TAI{Sec: seconds + 12*Year, Asec: nsec * Nanosecond}
}

// FromTime converts time t to TAI time, including handling of leap seconds
func FromTime(t time.Time) TAI {
	t = t.UTC()
	unix := t.Unix()
	nsec := t.Nanosecond()
	tai := Unix(unix, int64(nsec))
	skew := skew(t)
	tai.Sec -= skew
	return tai
}
