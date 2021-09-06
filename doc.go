/* Package tai provides functionality for International Atomic Time (TAI).

TAI times are unaffected by leapseconds and have an epoch of Jan 1 1958 00:00:00.

TAI times are comprised of one integer containing the number of whole seconds
since epoch and a second with the number of attoseconds of subsecond time elapsed.

TAI times only exist in the UTC time zone.

The resolution of this package exceeds most any clock, as of the end of year 2021.

TAI times support most of the same functions as a time.Time, based on the
Gregorian calendar.  They may be formatted similar to time, but format uses
more strftime-like specifiers.

To report accurate time in UTC from TAI, pkg tai must be periodically updated
to reflect Bulletin C by IERS.  Importing programs simply must track pkg tai as
it is updated.  Long-running programs (uptime > 6mo, say) must use the
RegisterLeapSecond function to remain current.

If the leap second table becomes stale, conversion to/from stdlib time.Time
will become inaccurate, by the number of missing leap seconds.  Tai only utilizes
the table at that interface.

## FAQ

1) Why would I want to use this?

If you deal with the TAI timekeeping system, there are not (as of late 2021) any
alternative packages for Go.  TAI time is continuous and never repeats, properties
that the Universal Time system used by most computers do not have.

2) Why not stdlib time?

An alternative implementation of this package would utilize the stdlib time
package and skew the values by the relevant number of leapseconds.  However,
stdlib time has a more finite range compared to the 292 billion years of pkg tai.

Additionally, support for non-UTC timezones is needless complexity, as TAI time
only exists in the UTC timezone.  As well, the nanosecond resolution of stdlib
time is restrictive for some applications interested in continuous time systems,
for which the attosecond resolution of this package is a tremendous improvement.

3) Is the package threadsafe?

Yes.  The leapsecond table is protected by a RWMutex.  This limits ultimate
concurrent performance when converting to/from stdlib Time values.

4) Why use global state for the leapsecond table?

The only significant benefit to a non-global leapsecond table is the ability
to use sharding to improve performance.  The performance of this package is
sufficient that additional speed would likely not be a significant gain to any
program.  A non-global leapsecond table must be duplicated N times over which
creates an enhanced opportunity for inconsistent state and incorrect programs.

5) Will there be more features to mimic a larger portion of the time package?

The time package is privileged with runtime hooks for some of its features, e.g.
timers.  The stdlib has a greater variety of serialization options for e.g. non-allocating approaches.
If desired, please implement them with a pull request.

The various methods for accessing parts of a time's representation such as the
weekday may be added at a later date.  It is unclear how valuable these are
when AsGreg exists.

6) How correct and bug free is this package?

There is a fairly complete set of unit tests that lead to the author's confidence
that the package works properly.  If you find a bug, please make an issue so it
can be fixed.

7) Why have a format syntax that is incompatible with Stdlib time?

Stdlib time is the ugly duckling in this respect.  pkg tai is more in keeping with
other similar tools.
*/
package tai
