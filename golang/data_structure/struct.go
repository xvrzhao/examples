package data_structure

import "fmt"

type Person struct {
	Name string
	Age  uint8
}

func RunStructAddressExample() {
	p := Person{Name: "Xavier", Age: 23}
	fmt.Println("The address of struct p is the same as the address of the first field of p:")
	fmt.Printf("address of p:\t\t%p\n", &p)
	fmt.Printf("address of p.Name:\t%p\n", &p.Name)
	fmt.Printf("address of p.Age:\t%p\n", &p.Age)
}
