package fibonacci

func fib(n int) int {
	if n <= 1 {
		return n
	}

	return fib(n-1) + fib(n-2)
}

func fibSequential(n int) int {
	if n <= 1 {
        return n
    }

    var n2, n1 int = 0, 1

    for i := 2 ; i < n ; i++ {
        n2, n1 = n1, n1 + n2
    }
 
	return n2 + n1
}

