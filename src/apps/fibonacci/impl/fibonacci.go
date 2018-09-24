package impl

type Fibonacci struct{}

func Fibo(n int) int {
	return F(n)
}

func F(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return F(n-1) + F(n-2)
	}
}
