package tai_test

import (
	"fmt"
	"testing"

	"github.com/brandondube/tai"
)

func TestIsLeapYearCorrect(t *testing.T) {
	cases := []struct {
		descr string
		inp   int
		exp   bool
	}{
		{"TestY1700", 1700, false},
		{"TestY1800", 1800, false},
		{"TestY1900", 1900, false},
		{"TestY2000", 2000, true},
		{"TestY2004", 2004, true},
		{"TestY0001", 0001, false},
		{"TestY0002", 0002, false},
		{"TestY0003", 0003, false},
		{"TestY0004", 0004, true},
	}
	for _, tc := range cases {
		t.Run(tc.descr, func(t *testing.T) {
			actual := tai.IsLeapYear(tc.inp)
			if actual != tc.exp {
				t.Fatalf("for year %d expected to get %v, got %v", tc.inp, tc.exp, actual)
			}
		})
	}
}

func TestZeroDayIsEpoch(t *testing.T) {
	y, m, d := tai.CivilFromDays(0)
	if y != 1958 {
		t.Fatal(fmt.Sprintf("day zero had year of %d, expected 1958", y))
	}
	if m != 1 {
		t.Fatal(fmt.Sprintf("day zero had month of %d, expected 1", m))
	}
	if d != 1 {
		t.Fatal(fmt.Sprintf("day zero had day of %d, expected 1", d))
	}
}

func TestCivilEpochIsZero(t *testing.T) {
	d := tai.DaysFromCivil(1958, 1, 1)
	if d != 0 {
		t.Fatal(fmt.Sprintf("epoch had day of %d, expected zero", d))
	}
}

func TestHammer(t *testing.T) {
	const (
		startYear = -4716
		endYear   = 10000
	)
	for y := startYear; y < endYear; y++ {
		for m := 1; m < 13; m++ {
			e := tai.DaysInMonth(m, y)
			for d := 1; d <= e; d++ {
				ta := tai.Date(y, m, d)
				g := ta.AsGreg()
				if g.Year != int64(y) || g.Month != uint8(m) || g.Day != uint8(d) {
					t.Fatal(fmt.Sprintf("input Y=%d, m=%d, d=%d failed, got Y=%d, m=%d, d=%d", y, m, d, g.Year, g.Month, g.Day))
				}
			}
		}
	}
}
