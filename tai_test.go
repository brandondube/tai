package tai_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/brandondube/tai"
)

func TestFuzzTaiToGreg(t *testing.T) {
	fuzzTaiToGreg(t, 100)
}
func fuzzTaiToGreg(t *testing.T, cases int) {
	for i := 0; i < cases; i++ {
		year := rand.Intn(2000)
		month := rand.Intn(12)
		day := rand.Intn(28) // min # of days in a month
		hour := rand.Intn(24)
		minute := rand.Intn(60)
		second := rand.Intn(60)
		tm := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
		unix := tm.Unix()
		ti := tai.Unix(unix, 0)
		date := ti.AsGreg()
		var failparts []string
		if date.Year != int64(year) {
			failparts = append(failparts, fmt.Sprintf("wrong year: got %d expected %d", date.Year, year))
		}
		if date.Month != uint8(month) {
			failparts = append(failparts, fmt.Sprintf("wrong month: got %d expected %d", date.Month, month))
		}
		if date.Day != uint8(day) {
			failparts = append(failparts, fmt.Sprintf("wrong day: got %d expected %d", date.Day, day))
		}
		if date.Hour != uint8(hour) {
			failparts = append(failparts, fmt.Sprintf("wrong hour: got %d expected %d", date.Hour, hour))
		}
		if date.Minute != uint8(minute) {
			failparts = append(failparts, fmt.Sprintf("wrong minute: got %d expected %d", date.Minute, minute))
		}
		if date.Sec != uint8(second) {
			failparts = append(failparts, fmt.Sprintf("wrong sec: got %d expected %d", date.Sec, second))
		}
		if date.Asec != 0 {
			failparts = append(failparts, fmt.Sprintf("wrong subsec: got %d expected %d", date.Asec, 0))
		}
		if len(failparts) != 0 {
			t.Fatal(strings.Join(failparts, "\n"))
		}
	}
}

func TestLessSpecialCasesGreg(t *testing.T) {
	cases := []struct {
		descr string
		inp   int
		exp   tai.Greg
	}{
		{"Positive1Y", +1 * tai.Year, tai.Greg{Year: 1959, Month: 1, Day: 1}}, // all others zero
		{"Positive2Y", +2 * tai.Year, tai.Greg{Year: 1960, Month: 1, Day: 1}}, // all others zero
		{"Negative1Y", -1 * tai.Year, tai.Greg{Year: 1957, Month: 1, Day: 1}}, // all others zero
		{"Negative2Y", -2 * tai.Year, tai.Greg{Year: 1956, Month: 1, Day: 1}}, // all others zero
		{"Negative3Y", -3 * tai.Year, tai.Greg{Year: 1955, Month: 1, Day: 1}}, // all others zero
		{"Negative4Y", -4 * tai.Year, tai.Greg{Year: 1954, Month: 1, Day: 1}}, // all others zero
		{"Positive1Y1M1D", 1*tai.Year + 32*tai.Day, tai.Greg{Year: 1959, Month: 2, Day: 2}},
	}
	for _, tc := range cases {
		t.Run(tc.descr, func(t *testing.T) {
			actual := tai.TAI{Sec: int64(tc.inp)}.AsGreg()
			if !actual.Eq(tc.exp) {
				t.Fatalf("expected %+v, got %+v", tc.exp, actual)
			}
		})
	}
}
func TestZeroTaiIsEpoch(t *testing.T) {
	var ta tai.TAI
	date := ta.AsGreg()
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
	if date.Minute != 0 {
		failparts = append(failparts, fmt.Sprintf("wrong minute: got %d expected %d", date.Minute, 0))
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

func TestUnixEpoch(t *testing.T) {
	var ta tai.TAI
	ta.Sec = 12 * tai.Year
	date := ta.AsGreg()
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
	if date.Minute != 0 {
		failparts = append(failparts, fmt.Sprintf("wrong minute: got %d expected %d", date.Minute, 0))
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

// top level result of these two benchmarks: can reduce space by > 50% without
// compromising time -> do so (keep changes)

func BenchmarkAsGregWithoutFmt(b *testing.B) {
	// 23.67 ns with all int64s
	// 22.99 ns with some uint8s
	now := tai.Now()
	for i := 0; i < b.N; i++ {
		now.AsGreg()
	}
}

func BenchmarkAsGregWithFmt(b *testing.B) {
	// 369.8 ns with all int64s
	// 364.6 ns with some uint8s
	now := tai.Now()
	for i := 0; i < b.N; i++ {
		g := now.AsGreg()
		fmt.Sprintf("%d %d %d %d %d %d %d", g.Year, g.Month, g.Day, g.Hour, g.Minute, g.Sec, g.Asec)
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
