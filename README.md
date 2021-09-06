# tai

Package *tai* provides support for [International Atomic Time](https://en.wikipedia.org/wiki/International_Atomic_Time).

## Usage

TAI values may be created directly,

```go
t := tai.TAI{Sec: 123456789, Asec: 300 * tai.Millisecond}
```

To break a TAI value into is constituent Calendar parts:

```go
g := t.AsGreg()
// g.Year
// g.Month
// g.Day
// g.Hour
// g.Minute
// g.Sec
```

"Greg" is shorthand for Gregorian, the calendar system used by most of the world.

## Formatting

```go
fmt.Println(tai.Now().Format(tai.RFC3339Micro)
// 2021-09-03T22:03:56.991894Z
```

## Stdlib compatibility

Convert to and from stdlib time values
```go
tai.FromTime(t.AsTime())
```

Compatible with the same UNIX notation as stdlib,
```go
tai.Unix(secs, nsecs).Unix() // back to sec/nsec
```

## More

See [pkg.go.dev](https://pkg.go.dev/github.com/brandondube/tai).

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
can be fixed.  The "Hammer" test ensures that all dates from the Julian epoch
(-4716) to the year 10,000 are understood correctly.  This is over 5.7 million
test cases.  An additional million fuzz test cases are run, as well as specific
tests for interesting or key moments.

7) Why have a format syntax that is incompatible with Stdlib time?

Stdlib time is the ugly duckling in this respect.  pkg tai is more in keeping with
other similar tools.
