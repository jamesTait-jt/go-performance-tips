package preallocate

import "testing"

func Benchmark_GrowSlice(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0 ; i < b.N ; i++ {
		growSlice()
	}
}

func Benchmark_GrowSlicePreAllocate(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0 ; i < b.N ; i++ {
		growSlicePreAllocate()
	}
}