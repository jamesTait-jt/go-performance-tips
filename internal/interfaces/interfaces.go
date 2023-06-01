package interfaces

type NoOpper interface {
	Int(n int)
	IntPtr(n *int)
}

type nothing struct {}

func (no nothing) Int(n int) {}

func (no nothing) IntPtr(n *int) {}

//go:noinline
func NoOpInt(no NoOpper, n int) {
	no.Int(n)
}

//go:noinline
func NoOpIntPtr(no NoOpper, n int) {
	no.IntPtr(&n)
}

