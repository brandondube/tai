/* Package tai provides functionality for measuring International Atomic Time (TAI).

TAI moments exclude leapseconds and have an epoch of Jan 1 1958 00:00:00.

TAI times are comprised of one integer containing the number of whole seconds
since epoch and a second with the number of attoseconds of subsecond time elapsed.

TAI times only exist in the UTC time zone.

NIST-F1 is the primary time and frequency standard of the United States and
contributes to the TAI standard.  Its uncertainty is 3.1e-16; approximately
two orders of magnitude larger than the representation in this package, 1e-18.

The resolution (accuracy) of this package exceeds most any clock, as of the end
of year 2021.

TAI times support most of the same functions as a time.Time, based on the
Gregorian calendar.  They may be formatted similar to time, but format uses
more C-like specifiers.

*/
package tai
