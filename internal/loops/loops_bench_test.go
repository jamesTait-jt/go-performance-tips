package loops

import "testing"

func BenchmarkForI(b *testing.B) {
	sl := []bigObj{
		{id: 0},
		{id: 1},
		{id: 2},
		{id: 3},
		{id: 4},
		{id: 5},
	}

	for i := 0 ; i < b.N ; i++ {
		forI(sl)
	}
}

func BenchmarkFoRange(b *testing.B) {
	sl := []bigObj{
		{id: 0},
		{id: 1},
		{id: 2},
		{id: 3},
		{id: 4},
		{id: 5},
	}

	for i := 0 ; i < b.N ; i++ {
		forRange(sl)
	}
}