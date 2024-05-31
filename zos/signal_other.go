//go:build windows

package zos

// Provide a fake "shim" for SIGUSR1 and SIGUSR2; adding support for listening
// for these usually isn't critical and an optional nice-to-have feature. It's
// okay if this won't do anything on Windows.
type fakeSignal int

func (fakeSignal) String() string { return "fake signal" }
func (fakeSignal) Signal()        {}

const (
	SIGUSR1 = fakeSignal(-1)
	SIGUSR2 = fakeSignal(-1)
)
