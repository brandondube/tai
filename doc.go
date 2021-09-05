/* Package tai provides functionality for International Atomic Time (TAI).

TAI times exclude leapseconds and have an epoch of Jan 1 1958 00:00:00.

TAI times are comprised of one integer containing the number of whole seconds
since epoch and a second with the number of attoseconds of subsecond time elapsed.

TAI times only exist in the UTC time zone.

The resolution of this package exceeds most any clock, as of the end of year 2021.

TAI times support most of the same functions as a time.Time, based on the
Gregorian calendar.  They may be formatted similar to time, but format uses
more strftime-like specifiers.

To report accurate time in UTC from TAI, programs that include pkg tai must
periodically check IERS bulletin C and register new leap seconds.  The package
includes variables that indicate when the leapsecond table was last updated and
when the leapsecond table expires.  An expired leap table may not contain new
leap seconds that have occurred, however issues of bulletin C are not guaranteed
to announce a new leap second.
*/
package tai
