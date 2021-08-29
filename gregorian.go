package tai

import (
	"fmt"
)

// maxint64 / seconds per year = 292277024626
// 292,277,024,626
// year 292 billion is when this becomes invalid
// (perfectly fine)
const (
	notAMonth = iota
	January
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)
const (
	Monday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

const (
	// to go from TAI to Unix, add the skew; subtract for Unix=>TAI
	unixSkewFwd = 12 * Year
	// AKA the UT2 instant
	epochJulianDay  = 2436204.5   // this is here to preserve the truth of the value
	epochJulianDayI = 2436204 + 1 // I = integer or truncated; use by offsetting 12h of seconds first
)

// IsLeapYear returns true if year is a leap year, false if
// year is not a leap year.
func IsLeapYear(year int) bool {
	/* per USNO,
	Every year that is exactly divisible by four is a leap year,
	except for years that are exactly divisible by 100,
	but these centurial years are leap years if they are exactly divisible by 400.
	For example, the years 1700, 1800, and 1900 are not leap years,
	but the years 1600 and 2000 are.
	*/
	if year < 1 {
		panic(fmt.Sprintf("tai/IsLeapYear: got year < 1 %d, not part of Gregorian Calendar", year))
	}
	if year%4 == 0 { // every year that is exactly divisible by four
		if year%100 == 0 { // except for years that are exactly divisible by 100
			if year%400 == 0 { // if they are exactly divisible by 400
				return true
			}
			return false
		}
		return true
	}
	return false
}

// SecsToJulianDay converts seconds since TAI epoch to Julian Day Number
func SecsToJulianDay(secs int64) int64 {
	secs += twelveHours
	return secs/Day + epochJulianDayI
}

// JulianDayToGregorianCalendar returns the day, month, and year in the Gregorian
// calendar corresponding to the given Julian Day Number
func JulianDayToGregorianCalendar(J int64) (Y, M, D int64) {
	// Algorithm devised by Edward Graham Richards, converted to Go by
	// Brandon Dube
	// see: https://en.wikipedia.org/wiki/Julian_day#Julian_or_Gregorian_calendar_from_Julian_day_number
	const (
		y = 4716
		j = 1401
		m = 2
		n = 12
		r = 4
		p = 1461
		v = 3
		u = 5
		s = 153
		w = 2
		B = 274277
		C = -38
	)
	f := J + j + (((4*J+B)/146097)*3)/4 + C
	e := r*f + v
	g := (e % p) / r
	h := u*g + w
	D = (h%s)/u + 1
	M = (h/s+m)%n + 1
	Y = (e / p) - y + (n+m-M)/n
	return // uncommon; named returns in Go do not require list of return "args"
}

// Greg represents a moment in the Proleptic Gregorian Calendar and the UTC time system
//
// Calculations that return Gregs must include leap seconds
type Greg struct {
	// formerly all i64 = 7x8 = 56 B
	// now two i64 + 5 u8 = 21 B
	Asec   int64
	Year   int64
	Month  uint8
	Day    uint8
	Hour   uint8
	Minute uint8
	Sec    uint8
}
