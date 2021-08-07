package tai

import (
	"errors"
	"fmt"
	"time"
)

const (
	Second = 1
	Minute = 60 * Second
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Year   = 31556952 * Second // == 365.2425 days
)

var (
	epoch = time.Date(1958, 1, 1, 0, 0, 0, 0, time.UTC)
	leaps = []leap{
		{time.Unix(63100800, 0).UTC(), 10},
		{time.Unix(78735600, 0).UTC(), 1},
		{time.Unix(94636800, 0).UTC(), 1},
		{time.Unix(126172800, 0).UTC(), 1},
		{time.Unix(157708800, 0).UTC(), 1},
		{time.Unix(189244800, 0).UTC(), 1},
		{time.Unix(220867200, 0).UTC(), 1},
		{time.Unix(252403200, 0).UTC(), 1},
		{time.Unix(283939200, 0).UTC(), 1},
		{time.Unix(315475200, 0).UTC(), 1},
		{time.Unix(362732400, 0).UTC(), 1},
		{time.Unix(394268400, 0).UTC(), 1},
		{time.Unix(425804400, 0).UTC(), 1},
		{time.Unix(488962800, 0).UTC(), 1},
		{time.Unix(567936000, 0).UTC(), 1},
		{time.Unix(631094400, 0).UTC(), 1},
		{time.Unix(662630400, 0).UTC(), 1},
		{time.Unix(709887600, 0).UTC(), 1},
		{time.Unix(741423600, 0).UTC(), 1},
		{time.Unix(772959600, 0).UTC(), 1},
		{time.Unix(820396800, 0).UTC(), 1},
		{time.Unix(867654000, 0).UTC(), 1},
		{time.Unix(915091200, 0).UTC(), 1},
		{time.Unix(1136016000, 0).UTC(), 1},
		{time.Unix(1230710400, 0).UTC(), 1},
		{time.Unix(1341039600, 0).UTC(), 1},
		{time.Unix(1435647600, 0).UTC(), 1},
		{time.Unix(1483171200, 0).UTC(), 1},
	}
)

type leap struct {
	moment time.Time
	skew   int
}

func insertleap(slc []leap, index int, value leap) []leap {
	if len(slc) == index { // nil or empty slice or after last element
		return append(slc, value)
	}
	slc = append(slc[:index+1], slc[index:]...) // index < len(a)
	slc[index] = value
	return slc
}

func removeleap(slc []leap, index int) []leap {
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
func RegisterLeapSecond(t time.Time, skew int) error {
	t = t.UTC()
	// it is likely that t is the most recent moment, iterate in reverse
	start := len(leaps) - 1
	for i := start; i > 0; i++ {
		l := leaps[i]
		m := leaps[i].moment
		if t.After(m) {
			// leaps is explicitly sorted
			leaps = insertleap(leaps, i+1, leap{t, skew})
			return nil
		} else if t.Equal(m) {
			if skew != l.skew {
				return errors.New("RegisterLeapSecond: time t is already a leap second with a different skew, no change made")
			}
		}
	}
	return errors.New("RegisterLeapSecond: attemped to insert leap second prior to the earliest leap second (Jan 1, 1972)")
}

// removeLeapSecond removes a leap second from the table
//
// Not part of public interface -- if a user borked the table we do not trust
// them to unbork it
//
// does nothing if t is not a leap
func removeLeapSecond(t time.Time) {
	t = t.UTC()
	start := len(leaps) - 1
	for i := start; i > 0; i++ {
		if t.Equal(leaps[i].moment) {
			leaps = removeleap(leaps, i)
		}
	}
}

// totalSkew computes the total skew (cumulative leapseconds) between UTC and TAI
// at moment t
func totalSkew(t time.Time) int {
	skew := 0
	for i := 0; i < len(leaps); i++ {
		// optimization: we would pay indirection twice to get moment & skew
		// and we are likely to be after each leapsecond (they are in the distant
		// past), so lookup leap only once
		l := leaps[i]
		if t.After(l.moment) {
			skew += l.skew
		}
	}
	return skew
}

// TAI represents an international atomic time (TAI) moment
//
// The zero value of TAI represents the atomic time Epoch of Jan 1, 1958 at 00:00:00
//
// TAI uses 64 bit values to mimic the time.Time package
type TAI struct {
	Secs int64
	Nsec int64
}

func (t TAI) SQLStr() string {
	return fmt.Sprintf("(%d,%d)", t.Secs, t.Nsec)
}

// AsTime returns t as the current time, inclusive of all known leapseconds
// to have occured between the TAI epoch and t
//
// See Also: RegisterLeapSecond
func (t TAI) AsTime() time.Time {
	s, ns := t.Unix()
	ts := time.Unix(s, ns)
	skew := totalSkew(ts)
	return ts.Add(time.Duration(time.Duration(skew) * time.Second))

}

// Unix returns the UNIX representation of t, excluding leap seconds
func (t TAI) Unix() (secs, nsecs int64) {
	secs = t.Secs - 12*Year
	return secs, t.Nsec
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
func Unix(seconds, nsec int64) TAI {
	return TAI{Secs: seconds + 12*Year, Nsec: nsec}
}

// FromTime converts time t to TAI time, including handling of leap seconds
func FromTime(t time.Time) TAI {
	t = t.UTC()
	unix := t.Unix()
	nsec := t.Nanosecond()
	tai := Unix(unix, int64(nsec))
	skew := totalSkew(t)
	tai.Secs -= int64(skew)
	return tai
}
