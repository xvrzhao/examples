package interfaces

func RunAssignStruct() {
	var (
		speaker Speaker
		walker  Walker
	)

	man := Man{}           // variable man implements the `Walk` method
	manPointer := new(Man) // variable manPointer implements both the `Speak` and the `Walk` method

	walker = man        // valid
	walker = manPointer // valid

	speaker = manPointer // valid
	// speaker = man	 // invalid, compiler error, variable man does not implement the `Speak` method

	use(speaker, walker)
}

// RunAssignInterface demonstrates that the success of the assignment between interfaces depends on
// whether the interface has implemented methods, not whether the underlying type of interface has implemented those methods.
//
// But, the type assertion between interfaces just depends on whether the underlying type of left has implemented the methods of right.
func RunAssignInterface() {
	var speaker Speaker
	speaker = new(Man)

	var speakWalker SpeakWalker
	//speakWalker = speaker // invalid, compiler error, speaker can not be assigned to speakWalker because speaker does not implement speakWalker, even though the underlying struct *Man implements that.

	speakWalker = new(Man)
	speaker = speakWalker // valid.

	use(speaker)
}
