package main

import (
	"adrpc"
	"fmt"
	"reflect"
	"sync"
	"time"
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

type S struct {
	k int
}

//第一个输入参数是传入结构，第二个参数是返回数指针
func (ss S) Add(a [2]int, b *int) error {
	*b = a[0] + a[1]
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

	m := sync.Map{}
	fmt.Println(m.Load("dasd"))
	m.LoadOrStore("d", 1)
	fmt.Println(m.Load("d"))
	m.Delete("d")
	fmt.Println(m.Load("d"))

	done := make(chan struct{}, 1)

	go func() {
		// 发送HTTP请求
		time.Sleep(1233 * time.Millisecond)
		done <- struct{}{}
	}()

	select {
	case <-done:
		fmt.Println("call successfully!!!")

	case <-time.After(time.Duration(800 * time.Millisecond)):
		fmt.Println("timeout!!!")
	}

	ss := adrpc.NewServer(1231)
	ss.Register(S{})
	go ss.StartServer(":1234")

	a := reflect.ValueOf([2]int{1, 16})

	var r = reflect.ValueOf(new(int))
	body := adrpc.Body{
		ServiceMethod: "S.Add",
		Args:          a,
		Reply:         r,
	}
	ss.Handle(&body)
	fmt.Println(body.Reply.Elem())

}
