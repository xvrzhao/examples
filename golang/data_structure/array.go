package data_structure

import "fmt"

func RunArrayAddressExample() {
	a := [3]int{1, 2, 3}
	fmt.Println("The address of the array a is the same as the address of the first element of a:")
	fmt.Printf("address of a:\t\t%p\n", &a)
	fmt.Printf("address of a[0]:\t%p\n", &a[0])
	fmt.Printf("address of a[1]:\t%p\n", &a[1])
	fmt.Printf("address of a[2]:\t%p\n", &a[2])

	s := a[:]
	fmt.Printf("The slice s points to:\t%p\n", s)
}
