package main

import (
	"adrpc"
	"adrpc/codec"
	"fmt"
	"log"
	"time"
)

type S struct {
	k int
}

type Args struct {
	Num1 int
	Num2 int
}

//第一个输入参数是传入结构，第二个参数是返回数指针
func (ss S) Add(a Args, b *int) error {
	*b = a.Num1 + a.Num2
	return nil
}

func StartServer() {
	server := adrpc.NewServer(12345)
	var s S
	server.Register(s)
	server.StartServer("127.0.0.1:1234")
}

func main() {
	fmt.Println("---------client test------------")
	go StartServer()
	client, err := adrpc.NewClient("127.0.0.1:1234", codec.JsonType, time.Microsecond*10000000)
	if err != nil {
		log.Fatalf("创建客户端失败")
	}
	var args = Args{
		Num1: 10,
		Num2: 2,
	}
	var reply = new(int)
	call := client.AddCall("S.Add", &args, reply)
	fmt.Println(call)
	client.Send(1, 12345, call)
	H, B := client.Receive()
	fmt.Print("接受完毕：")

	fmt.Println(*H, (*B).Reply)

}
