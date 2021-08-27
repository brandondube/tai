package tai

import "fmt"

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
	// indexed by month directly
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
)

const (
	// 1958:01:01
	epochYear  = 1958
	epochMonth = January
	epochDay   = 1
	epochSec   = 0
	epochAsec  = 0

	// to go from TAI to Unix, add the skew; subtract for Unix=>TAI
	unixSkewFwd = 12 * Year
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
