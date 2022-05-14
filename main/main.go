package main

import (
	"encoding/json"
	"fmt"
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
type Args struct {
	Num1 int
	Num2 int
}

//第一个输入参数是传入结构，第二个参数是返回数指针
func (ss S) Add(a Args, b *int) error {
	*b = a.Num1 + b.Num2
	return nil
}

type Monster struct {
	Name     string
	Age      int
	Birthday string
	Sal      float64
}

func main() {

	jsonStr := "{\"Name\":\"铁牛\", \"Age\":18,\"Birthday\":\"2020-02-02\",\"Sal\":1}"
	// 定义一个Monster实例
	var monster Monster
	err := json.Unmarshal([]byte(jsonStr), &monster)
	if err != nil {
		fmt.Printf("unarshar err=%v", err)
	}
	fmt.Printf("反序列化后 monster=%v\n", monster)

}
