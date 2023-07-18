package parse

import (
	"testing"
)

func Benchmark_Person(b *testing.B) {

	in := "joe bloggs 25 182.3"

	b.ReportAllocs()
	for i := 0 ; i < b.N ; i++ {
		Person(in)
	}
}
func Benchmark_PersonEfficient(b *testing.B) {
	s := "joe bloggs 25 182.3"

	b.ReportAllocs()
	for i := 0 ; i < b.N ; i++ {
		PersonEfficient(s)
	}
}

func Benchmark_PersonBytes(b *testing.B) {
	bs := []byte("joe\x1fbloggs\x1f25\x1f182.3")

	b.ReportAllocs()
	for i := 0 ; i < b.N ; i++ {
		PersonBytes(bs)
	}
}

func Benchmark_PersonBytesUnsafe(b *testing.B) {
	bs := []byte("joe\x1fbloggs\x1f25\x1f182.3")

	b.ReportAllocs()
	for i := 0 ; i < b.N ; i++ {
		PersonBytesUnsafe(bs)
	}
}