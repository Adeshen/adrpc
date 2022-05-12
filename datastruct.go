package adrpc

import (
	"adrpc/codec"
	"reflect"
	"sync"
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

//通信数据格式

type Header struct {
	ClientName  string
	Seq         int
	MagicNumber uint64
	//用于确定是否来自注册中心的允许,接受到报头是会携带服务端的魔数，以便确认
}

type Body struct {
	ServiceMethod string
	Args          reflect.Value
	Reply         reflect.Value
	err           error
}

//服务端数据结构

type service struct { //动态调用函数部分
	Name string
	//服务的名字    一般会    服务名.方法名
	Typ reflect.Type
	//结构体的类型反射
	Rcvr reflect.Value
	//结构体的值反射
	Method map[string]*methodType
}

type methodType struct { //动态调用函数
	method reflect.Method
	//函数的反射    调用结构method.Fun.call(所属结构体的反射，参数1的值反射，参数2的值反射....)
	ArgType reflect.Type
	//参数类型的反射
	ReplyType reflect.Type
	numCalls  uint64
}

type Server struct {
	Services map[string]*service
	mu       sync.Mutex

	ReqN uint64 //完成请求数

	MagicNumber uint64

	sending sync.Mutex
	//互斥量保证发送时不会冲突
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
