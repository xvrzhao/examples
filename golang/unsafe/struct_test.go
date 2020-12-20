package unsafe

import (
	"fmt"
	"testing"
)

func TestStructOffsetField(t *testing.T) {
	StructOffsetField()
}

func TestPerson2Bytes(t *testing.T) {
	p := Person{
		Name: "xavier",
		Age:  23,
	}

	b := Person2Bytes(&p)
	fmt.Println(b)

	pp := Bytes2Person(b)
	fmt.Println(pp.Name, pp.Age) // xavier 23

	pp.Name = "xvr"
	fmt.Println(p.Name) // xvr
}
