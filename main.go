package main

import (
	"fmt"
)

func f1() {
	var a *string
	println(*a)
}
func f2() {
	defer func() {
		if err, ok := recover().(error); ok {
			fmt.Printf("recovered: %s\n", err)
		}
	}()
	f1()
}
func main() {
	f2()
}
