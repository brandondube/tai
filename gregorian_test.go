package tai_test

import (
	"testing"

	"github.com/brandondube/tai"
)

func TestIsLeapYearValidYears(t *testing.T) {
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

func TestIsLeapYearPanicsForInvalidYears(t *testing.T) {
	cases := []struct {
		descr string
		inp   int
	}{
		{"TestY0", 0},
		{"TestY-1", -1},
	}
	for _, tc := range cases {
		t.Run(tc.descr, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil { // failed to panic
					t.Fatalf("for year %d, expected IsLeapYear to panic", tc.inp)
				}
			}()
			tai.IsLeapYear(tc.inp)
		})
	}
}
