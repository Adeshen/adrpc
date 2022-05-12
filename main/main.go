package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name  string
	Age   int
	Score float64
}

func (s Student) Info() {
	fmt.Println("Name =", s.Name, "Age =", s.Age, "Score =", s.Score)
}

func (s Student) Sfo() {
	fmt.Println("Name =", s.Name, "Age =", s.Age, "Score =", s.Score)
	return
}

type s struct {
	k int
}

func (ss s) Add(a int, b int) error {
	ss.k = a + b
	return nil
}

func main() {
	fmt.Println("嗨客网(www.haicoder.net)")
	var p = Student{
		Name:  "HaiCoder",
		Age:   10,
		Score: 99,
	}
	personValue := reflect.ValueOf(p)

	infoFunc := personValue.MethodByName("Info")
	infoFunc.Call([]reflect.Value{})

	fmt.Println(personValue.NumMethod())
	fmt.Println(personValue.Type().Method(0).Func.Call([]reflect.Value{personValue}))

}
