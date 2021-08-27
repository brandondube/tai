package tai_test

import (
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
