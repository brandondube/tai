package tai_test

import (
	"fmt"
	"testing"

	"github.com/brandondube/tai"
)

// 150 ns prev
// 3.048 ns post
// ~= 50 speedup
// => 18.22 ns with RWlock
// probably as good as can be done with safety

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
