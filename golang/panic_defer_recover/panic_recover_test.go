package panic_defer_recover

import "testing"

func TestReturnedInterfaceValue(t *testing.T) {
	ReturnedInterfaceValue()
}

func TestChildGoroutinePanic(t *testing.T) {
	ChildGoroutinePanic()
}

func TestPanicTerminateWholeProgram(t *testing.T) {
	PanicTerminateWholeProgram()
}

func TestRecoverMustInDefer(t *testing.T) {
	RecoverMustInDefer()
}
