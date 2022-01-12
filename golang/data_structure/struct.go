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

func RunStructPointerExample() {
	p1 := &Person{Name: "p1"}
	p2 := &Person{Name: "p2"}

	fmt.Printf("p1: %p %v\n", p1, p1)
	fmt.Printf("p2: %p %v\n", p2, p2)

	p1 = p2 // p1 所存指针改变，p1 与 p2 同时指向一块内存区域

	fmt.Printf("p1: %p %v\n", p1, p1)
	fmt.Printf("p2: %p %v\n", p2, p2)

	// ---

	p1 = &Person{Name: "p1"}
	p2 = &Person{Name: "p2"}

	fmt.Printf("p1: %p %v\n", p1, p1)
	fmt.Printf("p2: %p %v\n", p2, p2)

	*p1 = *p2 // p1 所存指针不变，仅仅将 p2 所指向的对象的值赋给 p1 所指向的对象

	fmt.Printf("p1: %p %v\n", p1, p1)
	fmt.Printf("p2: %p %v\n", p2, p2)
}
