package implrpc

import (
	"errors"
	"shared/shared"
)

type Calculator struct{}

func (t *Calculator) Mul(args *shared.Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Calculator) Sub(args *shared.Args, reply *int) error {
	*reply = args.A - args.B
	return nil
}
func (t *Calculator) Add(args *shared.Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}
func (t *Calculator) Div(args *shared.Args, quo *shared.Quotient) error {
	if args.B == 0 {
		return errors.New("Divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

