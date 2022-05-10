package adrpc

import (
	"net"
	"reflect"
)

//客户端数据结构
type Call struct {
	method string
	args   interface{}
	reply  interface{}
	seq    int
}

type Client struct {
	conn net.Conn
	seq  int
	call chan *Call
}

//-----------

//负载均衡客户端管理
type ClientMaster struct {
	client map[string]*Client
}

//-----------

//服务端数据结构
type Request struct {
	method string
	args   interface{}
	reply  interface{}
	seq    int
}

type Service struct {
	name    string
	methods map[string]*Method
}

type Method struct {
	name string
	fun  reflect.Value
}

type Server struct {
	services map[string]*Service
	conn     net.Conn
}

///----------

//注册中心+服务器负载分配
type Register struct {
	servers    map[string]string //service.method -> addr
	serverload map[string]int    // addr -> load
}
