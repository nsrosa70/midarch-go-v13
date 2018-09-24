package implrpc

import (
	"shared/shared"
)

type Fibonacci struct{}

func (f Fibonacci) Fibo(n shared.Args, reply *int) error {
	*reply = F(n.A)

	return nil
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


