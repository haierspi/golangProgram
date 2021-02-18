package main

import (
	"fmt"
	"reflect"
)

type person struct {
	name string
}

type student struct {
	*person
	number string
}

func main() {
	var p = new(student)
	p.number = "20201201"
	p.person = new(person)
	p.name = "guanglai"
	fmt.Print(reflect.TypeOf(p), p.number, p.name)
}
