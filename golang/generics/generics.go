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

// Tree is a generic type.
type Tree[T string | int] struct {
	left, right *Tree[T]
	value       T
}

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
	// Providing type arguments to a generic function or a generic type, then we get
	// a non-generic one, this process is called `instantiation`.
	//
	// Note that any generic function or generic type must be explicitly or implicitly
	// instantiated before it can be used.

	// example 1:
	// Explicit instantiation of generic function.
	handleMan := handleHuman[*Man]
	handleMan(new(Man))
	// handleMan(new(Woman)) // bad

	// example 2:
	// Implicit instantiation of generic function, omit the type argument because the compiler can infer them.
	handleHuman(new(Man))

	// example 2:
	// Generic type must be explicitly instantiated before it can be used.
	var intTree Tree[int]
	_ = intTree
	// or
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

func ExampleOfTypeParameter() {
	// Type parameter list can only be declared with outermost function or type,
	// and can't be declared with closure or anonymous function.

	// error example 1:
	//   f := func[T int | float64] (n T) {} // syntax error: function type must have no type parameters

	// error example 2:
	//   func receiveClosure(closure func[T int | float64](a, b T)) { ... } // syntax error: function type must have no type parameters

	// correct example:
	// func receiveClosure[T int | float64](closure func(a, b T)) {}
}

func receiveClosure[T int | float64](closure func(a, b T)) {}

func ExampleOfGenericFuncWithClosureParam() {
	// First instantiate and then call, the type in the closure parameter must be repalced with the outer type argument.
	receiveClosure[int](func(a, b int) {})
	receiveClosure[float64](func(a, b float64) {})

	// or omit the type argument
	receiveClosure(func(a, b int) {})
	receiveClosure(func(a, b float64) {})
}

type number interface {
	~int | ~float64
}

type myint int
type myfloat64 float64

func test[T number](n T) {}

func ttt() {
	test(23)
	test(1.2)
	test(myint(1))
	test(myfloat64(1))
}
