package race_and_sync

import "testing"

func TestIntRace1(t *testing.T) {
	IntRace1()
}

func TestIntRace2(t *testing.T) {
	IntRace2()
}

func TestIntAtomic(t *testing.T) {
	IntAtomic()
}

func TestStringRace(t *testing.T) {
	StringRace()
}

func TestBoolRace(t *testing.T) {
	BoolRace()
}
