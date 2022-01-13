package oop

import "fmt"

type father struct{}

func (f *father) eat() {
	fmt.Println("father eats")
}

func (f *father) live() {
	f.eat()
}

type son struct {
	*father
}

func (s *son) eat() {
	fmt.Println("son eats")
}

func RunExampleOfInherit() {
	son{new(father)}.live() // father eats, not son
}
