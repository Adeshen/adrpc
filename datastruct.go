package adrpc

import (
	"net"
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
	Conn     net.Conn //net.Dial() 创建连接
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
	Args   interface{}
	Reply  interface{}
	Seq    int
}

type Service struct {
	Name    string
	Methods map[string]*Method
}

type Method struct {
	Name string
	Fun  reflect.Value
}

type Server struct {
	Services map[string]*Service
	Conn     net.Conn
}

///----------

//注册中心+服务器负载分配
type Register struct {
	servers    map[string]string //service.method -> addr
	serverload map[string]int    // addr -> load value
}
