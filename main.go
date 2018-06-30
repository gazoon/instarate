package main

import "fmt"

type A interface {
	Foo() string
}
type B struct {
}

func (self *B) Foo() string {
	return ""
}
func main() {
	var a A = &B{}
	fmt.Printf("%T\nllh", a)
}
