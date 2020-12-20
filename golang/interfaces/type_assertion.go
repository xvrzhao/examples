package interfaces

import (
	"fmt"
	"log"
)

// RunTypeAssertionExample1 demonstrates how to determine the underlying struct type of a interface value.
func RunTypeAssertionExample1() {
	human := newHuman("Xavier")

	/*
	 * unchecked type assertion
	 */

	// girl := human.(*Girl) // panic as human is not *Girl
	// man := human.(*Man) // man, the value of type assertion, is type *Man

	/*
	 * checked type assertion, which is the best practice
	 */

	girl, ok := human.(*Girl) // not panic although human is not *Girl, girl is type *Girl and value is `zero value` of *Girl, ok is false
	if ok {
		fmt.Println(girl.Gender)
	} else {
		fmt.Println("human is not *Girl")
	}

	// man is type *Man that unwrapped from human interface, ok is true
	if man, ok := human.(*Man); ok {
		fmt.Println(man.Name)
	} else {
		fmt.Println("human is not *Man")
	}

	// just to judge
	if _, ok = human.(*Girl); ok {
		// ...
	}
}

// RunTypeAssertionExample2 demonstrates that how to determine if the underlying struct of a interface value implements other interfaces.
func RunTypeAssertionExample2() {
	human := newHuman("Xavier")
	if animal, ok := human.(Animal); ok { // the left of type assertion must be interface type
		animal.Walk()
	} else {
		log.Println("human is not Animal")
	}
}
