package tai_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/brandondube/tai"
)

func TestRegisterLeapSecondFunctions(t *testing.T) {
	err := tai.RegisterLeapSecond(1e12, 40)
	if err != nil {
		t.Fatal("non-nil err registering a leapsecond in the distant future", err)
	}
	tai.RemoveLeapSecond(1e12) // cleanup
}

func TestFuzzTaiToGreg(t *testing.T) {
	fuzzTaiToGreg(t, 1e6)
}
func fuzzTaiToGreg(t *testing.T, cases int) {
	for i := 0; i < cases; i++ {
		year := rand.Intn(2000)
		month := rand.Intn(12)
		day := rand.Intn(28) // min # of days in a month
		hour := rand.Intn(24)
		minute := rand.Intn(60)
		second := rand.Intn(60)
		if month == 0 {
			month = 1
		}
		if day == 0 {
			day = 1
		}
		ta := tai.Date(year, month, day).AddHMS(hour, minute, second)
		date := ta.AsGregorian()
		var failparts []string
		if date.Year != year {
			failparts = append(failparts, fmt.Sprintf("wrong year: got %d expected %d", date.Year, year))
		}
		if date.Month != month {
			failparts = append(failparts, fmt.Sprintf("wrong month: got %d expected %d", date.Month, month))
		}
		if date.Day != day {
			failparts = append(failparts, fmt.Sprintf("wrong day: got %d expected %d", date.Day, day))
		}
		if date.Hour != hour {
			failparts = append(failparts, fmt.Sprintf("wrong hour: got %d expected %d", date.Hour, hour))
		}
		if date.Min != minute {
			failparts = append(failparts, fmt.Sprintf("wrong minute: got %d expected %d", date.Min, minute))
		}
		if date.Sec != second {
			failparts = append(failparts, fmt.Sprintf("wrong sec: got %d expected %d", date.Sec, second))
		}
		if date.Asec != 0 {
			failparts = append(failparts, fmt.Sprintf("wrong subsec: got %d expected %d", date.Asec, 0))
		}
		if len(failparts) != 0 {
			failparts = append(failparts, fmt.Sprintf("input Year=%d, Month=%d, Day=%d, Hour=%d, Min=%d, Sec=%d", year, month, day, hour, minute, second))
			t.Fatal(strings.Join(failparts, "\n"))
		}
	}
}

func TestLessSpecialCasesGreg(t *testing.T) {
	cases := []struct {
		descr string
		inp   tai.TAI
		exp   tai.Gregorian
	}{
		{"Positive1Y", tai.Date(1959, 1, 1), tai.Gregorian{Year: 1959, Month: 1, Day: 1}}, // all others zero
		{"Positive2Y", tai.Date(1960, 1, 1), tai.Gregorian{Year: 1960, Month: 1, Day: 1}}, // all others zero
		{"Negative1Y", tai.Date(1957, 1, 1), tai.Gregorian{Year: 1957, Month: 1, Day: 1}}, // all others zero
		{"Negative2Y", tai.Date(1956, 1, 1), tai.Gregorian{Year: 1956, Month: 1, Day: 1}}, // all others zero
		{"Negative3Y", tai.Date(1955, 1, 1), tai.Gregorian{Year: 1955, Month: 1, Day: 1}}, // all others zero
		{"Negative4Y", tai.Date(1954, 1, 1), tai.Gregorian{Year: 1954, Month: 1, Day: 1}}, // all others zero
		{"Positive1Y1M1D", tai.Date(1959, 2, 2), tai.Gregorian{Year: 1959, Month: 2, Day: 2}},
		{"Negative1Y1M1D", tai.Date(1956, 2, 2), tai.Gregorian{Year: 1956, Month: 2, Day: 2}},
		{"DayOfChangeToGregorian", tai.Date(1582, tai.October, 15), tai.Gregorian{Year: 1582, Month: 10, Day: 15}},
		{"LastJulianDay", tai.Date(1582, tai.October, 4), tai.Gregorian{Year: 1582, Month: 10, Day: 4}},
		{"BrokenFuzzCase1NoHMS", tai.Date(81, 3, 15), tai.Gregorian{Year: 81, Month: 3, Day: 15}},
		{"BrokenFuzzCase1", tai.Date(81, 3, 15).AddHMS(11, 1, 18), tai.Gregorian{Year: 81, Month: 3, Day: 15, Hour: 11, Min: 1, Sec: 18}},
	}
	for _, tc := range cases {
		t.Run(tc.descr, func(t *testing.T) {
			actual := tc.inp.AsGregorian()
			if !actual.Eq(tc.exp) {
				t.Fatalf("expected %+v, got %+v", tc.exp, actual)
			}
		})
	}
}
func TestZeroTaiIsEpoch(t *testing.T) {
	var ta tai.TAI
	date := ta.AsGregorian()
	var failparts []string
	if date.Year != 1958 {
		failparts = append(failparts, fmt.Sprintf("wrong year: got %d expected %d", date.Year, 1958))
	}
	if date.Month != tai.January {
		failparts = append(failparts, fmt.Sprintf("wrong month: got %d expected %d", date.Month, tai.January))
	}
	if date.Day != 1 {
		failparts = append(failparts, fmt.Sprintf("wrong day: got %d expected %d", date.Day, 1))
	}
	if date.Hour != 0 {
		failparts = append(failparts, fmt.Sprintf("wrong hour: got %d expected %d", date.Hour, 0))
	}
	if date.Min != 0 {
		failparts = append(failparts, fmt.Sprintf("wrong minute: got %d expected %d", date.Min, 0))
	}
	if date.Sec != 0 {
		failparts = append(failparts, fmt.Sprintf("wrong sec: got %d expected %d", date.Sec, 0))
	}
	if date.Asec != 0 {
		failparts = append(failparts, fmt.Sprintf("wrong subsec: got %d expected %d", date.Asec, 0))
	}
	if len(failparts) != 0 {
		t.Fatal(strings.Join(failparts, "\n"))
	}
}

func TestTaiFormat(t *testing.T) {
	//2009-11-10 23:00:00
	ta := tai.Date(2009, 11, 10).AddHMS(23, 0, 0)
	out := ta.Format(tai.RFC3339)
	if out != "2009-11-10T23:00:00Z" {
		t.Fail()
	}
	out = ta.Format(tai.RFC3339Micro)
	if out != "2009-11-10T23:00:00.000000Z" {
		t.Fail()
	}

	out = ta.Format(tai.RFC3339Nano)
	if out != "2009-11-10T23:00:00.000000000Z" {
		t.Fail()
	}
}

func TestNowAsTimeEq(t *testing.T) {
	now := tai.Now()
	nowT := now.AsTime()
	nowT2 := time.Now()
	diff := nowT2.Sub(nowT)
	if diff < 0 {
		diff = -diff
	}
	if diff > 100*time.Millisecond {
		t.Fatal("tai now and stdlib now differ by > 100 msec")
	}
}

func TestFromTimeAsTimeRoundTrip(t *testing.T) {
	now := time.Now()
	now2 := tai.FromTime(now).AsTime()
	if !now.Equal(now2) {
		t.Fatal()
	}
}

func TestUnixEpoch(t *testing.T) {
	ta := tai.Tai(4383*tai.Day, 0)
	date := ta.AsGregorian()
	var failparts []string
	if date.Year != 1970 {
		failparts = append(failparts, fmt.Sprintf("wrong year: got %d expected %d", date.Year, 1970))
	}
	if date.Month != tai.January {
		failparts = append(failparts, fmt.Sprintf("wrong month: got %d expected %d", date.Month, tai.January))
	}
	if date.Day != 1 {
		failparts = append(failparts, fmt.Sprintf("wrong day: got %d expected %d", date.Day, 1))
	}
	if date.Hour != 0 {
		failparts = append(failparts, fmt.Sprintf("wrong hour: got %d expected %d", date.Hour, 0))
	}
	if date.Min != 0 {
		failparts = append(failparts, fmt.Sprintf("wrong minute: got %d expected %d", date.Min, 0))
	}
	if date.Sec != 0 {
		failparts = append(failparts, fmt.Sprintf("wrong sec: got %d expected %d", date.Sec, 0))
	}
	if date.Asec != 0 {
		failparts = append(failparts, fmt.Sprintf("wrong subsec: got %d expected %d", date.Asec, 0))
	}
	if len(failparts) != 0 {
		t.Fatal(strings.Join(failparts, "\n"))
	}
}

func BenchmarkTaiAsTime(b *testing.B) {
	now := tai.Now()
	for i := 0; i < b.N; i++ {
		now.AsTime()
	}
}

func BenchmarkTaiFormat(b *testing.B) {
	now := tai.Now()
	for i := 0; i < b.N; i++ {
		now.Format(tai.RFC3339Micro)
	}
}

// top level result of these two benchmarks: can reduce space by > 50% without
// compromising time -> do so (keep changes)

func BenchmarkAsGregorianWithoutFmt(b *testing.B) {
	// 23.67 ns with all int64s
	// 22.99 ns with some uint8s
	now := tai.Now()
	for i := 0; i < b.N; i++ {
		now.AsGregorian()
	}
}

func BenchmarkAsGregorianWithFmt(b *testing.B) {
	// 369.8 ns with all int64s
	// 364.6 ns with some uint8s
	now := tai.Now()
	for i := 0; i < b.N; i++ {
		g := now.AsGregorian()
		fmt.Sprintf("%d %d %d %d %d %d %d", g.Year, g.Month, g.Day, g.Hour, g.Min, g.Sec, g.Asec)
	}
}

func BenchmarkTimeWithoutFmt(b *testing.B) {
	// 35.92 ns; tai ~= 33% faster
	now := time.Now()
	for i := 0; i < b.N; i++ {
		now.Date()
		now.Second()
		now.Nanosecond()
	}
}

func BenchmarkTimeFormat(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		now.Format(time.RFC3339Nano)
	}
}
