package adrpc

import (
	"adrpc/codec"
	"reflect"
)

//客户端数据结构
type Call struct {
	Method string
	Args   interface{}
	Reply  interface{} //需要填写答案时使用
	Seq    int
}

type Client struct {
	// Conn     net.Conn //net.Dial() 创建连接
	NetIO    codec.Codec
	Seq      int
	Callchan chan *Call
}

//-----------

//负载均衡客户端管理
type ClientMaster struct {
	ClientMap map[string]*Client
}

//-----------

//服务端数据结构
type Request struct {
	Method string
	Args   reflect.Value
	Reply  reflect.Value
	Seq    int
}

type service struct { //动态调用函数部分
	Name   string
	Typ    reflect.Type
	Rcvr   reflect.Value
	Method map[string]*methodType
}

type methodType struct { //动态调用函数
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	numCalls  uint64
}

type Server struct {
	Services map[string]*service
}

///----------

//注册中心+服务器负载分配
type Register struct {
	servers    map[string]string //service.method -> addr
	serverload map[string]int    // addr -> load value
}

//---------------

//编码器

//采用json、god两种编码模式
