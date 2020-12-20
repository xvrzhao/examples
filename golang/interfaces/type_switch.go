package interfaces

import "fmt"

func RunTypeSwitchExample() {
	human := newHuman("Xavier")

	switch t := human.(type) {
	case *Man:
		fmt.Println("*Man", t.Name)
	case *Girl:
		fmt.Println("*Girl", t.Gender)
	default:
		fmt.Println("human is not *Man or *Girl")
	}

	// fmt.Println(t) // can not use t outside the switch scope

	switch human.(type) {
	case *Man:
		// todo
	case *Girl:
		// todo
	default:
		// todo
	}
}
