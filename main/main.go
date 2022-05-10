package main

import (
	"adrpc"
	"fmt"
	"reflect"
)

type A struct {
	A int
}

type Student struct {
	Name string
	Age  int
	Sex  string
}

func main() {
	var stu Student
	v := reflect.ValueOf(&stu)
	v.Elem().FieldByName("Name").SetString("caigy")
	v.Elem().FieldByName("Age").SetInt(18)
	v.Elem().FieldByName("Sex").SetString("male")

	fmt.Printf("Student: %+v", stu)

	var aa A

	fmt.Println(aa)

	aav := reflect.ValueOf(&aa)
	aav.Elem().FieldByName("A").SetInt(100)
	fmt.Println(aa)

	call := adrpc.Call{}
	fmt.Println(call)
	callv := reflect.ValueOf(&call)

	callv.Elem().FieldByName("reply").SetInt(10)
	fmt.Println(call)
}
