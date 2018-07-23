package main

import "fmt"

var c = 2

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
	fmt.Println(c)
}
