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

var (
	daysPerNonLeapMonth = [...]int{
		0,
		31,
		28,
		31,
		30,
		31,
		30,
		31,
		31,
		30,
		31,
		30,
		31,
	}
	daysPerLeapMonth = [...]int{
		0,
		31,
		29,
		31,
		30,
		31,
		30,
		31,
		31,
		30,
		31,
		30,
		31,
	}
	daysBeforeNonLeapMonth = [...]int{
		0,       // not a month
		0,       // January
		31,      // February
		31 + 28, // ...
		31 + 28 + 31,
		31 + 28 + 31 + 30,
		31 + 28 + 31 + 30 + 31,
		31 + 28 + 31 + 30 + 31 + 30,
		31 + 28 + 31 + 30 + 31 + 30 + 31,
		31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
		31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
		31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
		31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
		31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
	}
	monthNamesFull = [...]string{
		"not a month",
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}
	monthNamesAbbrev = [...]string{
		"NaM",
		"Jan",
		"Feb",
		"Mar",
		"Apr",
		"May",
		"Jun",
		"Jul",
		"Aug",
		"Sept",
		"Oct",
		"Nov",
		"Dec",
	}
	weekdayNames = [...]string{
		"Sunday",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
	}
	weekdayNamesAbbrev = [...]string{
		"Sun",
		"Mon",
		"Tue",
		"Wed",
		"Thu",
		"Fri",
		"Sat",
	}
)

const (
	eraYears   = 400
	eraYearsm1 = eraYears - 1
	epochDays  = 719468 - 4383 // 719468 == Jan 1 1970 from 0000 Mar 1
	yearDays   = 365
	eraDays    = 146097
	eraDaysm1  = eraDays - 1
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
			return year%400 == 0 // if they are exactly divisible by 400
		}
		return true
	}
	return false
}

/*
The functions:

DaysFromCivil
CivilFromDays
WeekdayFromDays
WeekdayDifference
NextWeekday
PrevWeekday
DaysFromSecsEpoch
SecsEpochFromDays

are adapted from Howard Hinnant's public domain algorithms

https://howardhinnant.github.io/date_algorithms.html

Thank you, Howard!
*/

// DaysFromCivil returns the number of days in the Gregorian calendar since
// Jan 1, 1958 from a year, month, and day
func DaysFromCivil(y, m, d int) int {
	if m <= 2 {
		y--
	}
	var era, doy int
	if y >= 0 {
		era = y / eraYears
	} else {
		era = (y - eraYearsm1) / eraYears
	}
	yoe := y - era*eraYears
	if m > 2 {
		m -= 3
	} else {
		m += 9
	}
	doy = (153*m+2)/5 + d - 1
	doe := yoe*yearDays + yoe/4 - yoe/100 + doy
	return era*eraDays + doe - epochDays
}

func CivilFromDays(days int) (y, m, d int) {
	days += epochDays
	var era, doe, yoe int
	if days >= 0 {
		era = days
	} else {
		era = days - eraDaysm1
	}
	era /= eraDays
	doe = days - era*eraDays
	yoe = (doe - doe/1460 + doe/36524 - doe/146096) / 365
	y = yoe + era*eraYears
	doy := doe - (365*yoe + yoe/4 - yoe/100)
	mp := (5*doy + 2) / 153
	d = doy - (153*mp+2)/5 + 1
	if mp < 10 {
		m = mp + 3
	} else {
		m = mp - 9
	}
	if m <= 2 {
		y++
	}
	return
}

// 0 == sunday, 6 == sat; not ISO (0 == monday)
func WeekdayFromDays(days int) int {
	if days >= -4 {
		return (days + 4) % 7
	}
	return (days+5)%7 + 6
}

func WeekdayDifference(d1, d2 int) int {
	d1 -= d2
	if d1 <= 6 {
		return d1
	}
	return d1 + 7
}

func NextWeekday(wd int) int {
	if wd < 6 {
		return wd + 1
	}
	return 0
}

func PrevWeekday(wd int) int {
	if wd > 0 {
		return wd - 1
	}
	return 6
}

func DaysFromSecsEpoch(secs int) int {
	return secs / Day
}

func SecsEpochFromDays(days int) int {
	return days * Day
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

// Before returns true if g is before o
func (g Greg) Before(o Greg) bool {
	t1 := FromGreg(g)
	t2 := FromGreg(o)
	return t1.Before(t2)
}

// After returns true if g is after o
func (g Greg) After(o Greg) bool {
	t1 := FromGreg(g)
	t2 := FromGreg(o)
	return t1.After(t2)
}

// Eq returns true if g and o represent the same instant in time
func (g Greg) Eq(o Greg) bool {
	return (g.Asec == o.Asec &&
		g.Year == o.Year &&
		g.Month == o.Month &&
		g.Day == o.Day &&
		g.Hour == o.Hour &&
		g.Minute == o.Minute &&
		g.Sec == o.Sec)
}
