package fibonacci

import "testing"

func Benchmark_Fib(b *testing.B) {

	scenarios := []struct{
		name string
		input int
	}{
		{
			name: "base_0",
			input: 0,
		}, {
			name: "base_1",
			input: 1,
		}, {
			name: "recursive",
			input: 10,
		}, {
			name: "recursive_big",
			input: 25,
		},
	}

	for _, bs := range scenarios {
		b.Run(bs.name, func(b *testing.B) {
			for i := 0 ; i < b.N ; i++ {
				fib(bs.input)
			}
		})
	}
}

func Benchmark_FibSequential(b *testing.B) {

	scenarios := []struct{
		name string
		input int
	}{
		{
			name: "base_0",
			input: 0,
		}, {
			name: "base_1",
			input: 1,
		}, {
			name: "recursive",
			input: 10,
		}, {
			name: "recursive_big",
			input: 25,
		},
	}

	for _, bs := range scenarios {
		b.Run(bs.name, func(b *testing.B) {
			for i := 0 ; i < b.N ; i++ {
				fibSequential(bs.input)
			}
		})
	}
}