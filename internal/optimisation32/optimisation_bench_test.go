package optimisation32

import "testing"

func BenchmarkString(b *testing.B) {
	b.Run("LessThan32", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0 ; i < b.N ; i++ {
			in := make([]byte, 31)
			_ = string(in)
		}
	})
	
	b.Run("EqualTo32", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0 ; i < b.N ; i++ {
			in := make([]byte, 32)
			_ = string(in)
		}
	})
	
	b.Run("GreaterThan32", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0 ; i < b.N ; i++ {
			in := make([]byte, 33)
			_ = string(in)
		}
	})
}