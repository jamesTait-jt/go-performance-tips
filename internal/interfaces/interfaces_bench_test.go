package interfaces

import (
	"testing"
)

func BenchmarkNoOpInt(b *testing.B) {
	b.ReportAllocs()

	no := nothing{}
	n := 10

	for i := 0 ; i < b.N ; i++ {
		NoOpInt(no, n)
	}
}

func BenchmarkNoOpIntPtr(b *testing.B) {
	b.ReportAllocs()

	no := nothing{}
	n := 10

	for i := 0 ; i < b.N ; i++ {
		NoOpIntPtr(no, n)
	}
}