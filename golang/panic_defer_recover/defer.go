package panic_defer_recover

import (
	"fmt"
	"os"
)

// DeferTrap prints results of all traps about defer.
// Run the command:
//   $ go test -v -run=DeferTrap -count=1 ./panic_defer_recover
func DeferTrap() {
	deferTrap1()
	fmt.Println("---")

	deferTrap2()
	fmt.Println("---")

	deferTrap3()
	fmt.Println("---")

	deferTrap4()
	fmt.Println("---")

	deferTrap5()
	fmt.Println("---")

	deferTrap6()()
	fmt.Println("---")

	fmt.Println(deferTrap7())
	fmt.Println("---")

	fmt.Println(deferTrap8())
	fmt.Println("---")

	fmt.Println(deferTrap9())
	fmt.Println("---")

	deferTrap10()
}

func deferTrap1() {
	fmt.Println("trap1:")
	x := 1
	y := 2
	defer calc("A", x, calc("B", x, y))
	x = 3
	defer calc("C", x, calc("D", x, y))
	y = 4
}

// LIFO

func deferTrap2() {
	defer fmt.Print("Xavier Zhao.\n")
	defer fmt.Print("is short for ")
	defer fmt.Print("xvrzhao ")
	fmt.Println("trap2:")
}

// closure

func deferTrap3() {
	a, b := 1, 2
	defer fmt.Println(a + b)
	a = 2
	fmt.Println("trap3:")
}

func deferTrap4() {
	a, b := 1, 2
	defer func() {
		fmt.Println(a + b)
	}()
	a = 2
	fmt.Println("trap4:")
}

func deferTrap5() {
	a, b := 1, 2
	defer func(a, b int) {
		fmt.Println(a + b)
	}(a, b)
	a = 2
	fmt.Println("trap5:")
}

func deferTrap6() func() {
	fmt.Println("trap6:")
	a := 1
	defer func() {
		a++
	}()
	return func() {
		fmt.Println(a)
	}
}

// return

func deferTrap7() int {
	fmt.Println("trap7:")
	a := 1
	defer func() {
		a++
	}()
	return a
}

func deferTrap8() (a int) {
	fmt.Println("trap8:")
	defer func() {
		a++
	}()
	return 1
}

func deferTrap9() (a int) {
	fmt.Println("trap9:")
	defer func(a int) {
		a++
	}(a)
	return 1
}

// exit

func deferTrap10() {
	defer fmt.Println("won't be executed")
	fmt.Println("trap10:")
	os.Exit(0)
}

func calc(i string, a, b int) int {
	sum := a + b
	fmt.Println(i, a, b, sum)
	return sum
}
