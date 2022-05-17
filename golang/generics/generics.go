package generics

type Man struct{}

func (Man) Talk()  {}
func (Man) Stand() {}

type Woman struct{}

func (Woman) Talk() {}

type Transgender struct{}

func (Transgender) Talk()  {}
func (Transgender) Stand() {}

type Human interface {
	Talk()
	Stand()
	*Man | *Woman
}

// handleHuman is a generic function.
func handleHuman[V Human](human V) {}

func ExampleOfTypeConstraint() {
	// If an interface contains type constraints, it can't been used to declare
	// an variable or a function parameter.
	//
	//   var human Human				// can't be compiled
	//   func handleHuman(human Human)	// can't be compiled

	// When an interface contains not only type constraints but also methods,
	// and is used to declare type parameters, then the corresponding type argument
	// must be one of these types and implement all of these methods.
	handleHuman(new(Man))
	// handleHuman(new(Woman))       // *Woman does not implement Human (missing method Stand)
	// handleHuman(new(Transgender)) // *Transgender does not implement Human
}

func ExampleOfInstantiation() {
	// Providing the type argument to the generic function, then we have a non-generic
	// one, this process is called `instantiation`.
	handleMan := handleHuman[*Man]
	handleMan(new(Man))
	// handleMan(new(Woman)) // bad
}

type Person[T *Man | *Woman] struct {
	gender string
	mate   T
}

type Tree[T string | int] struct {
	left, right *Tree[T]
	value       T
}

func ExampleOfGenericType() {
	// When using generic types, be aware that instantiation is required.

	_ = Tree[int]{
		left:  &Tree[int]{},
		right: nil,
		value: 0,
	}

	// or

	type stringTree = Tree[string]
	_ = stringTree{
		left:  nil,
		right: nil,
		value: "",
	}
}
