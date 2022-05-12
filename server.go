package adrpc

import (
	"adrpc/codec"
	"errors"
	"log"
	"net"
	"strings"
)

//以下方法按时间顺序排列，创建---注册服务------
///接受连接-------读取到来信息------处理请求-------调用服务------发送回信

func NewServer(magicnumber uint64) *Server {
	server := Server{
		MagicNumber: uint64(magicnumber),
		Services:    make(map[string]*service),
	}
	return &server
}

func (server *Server) StartServer(port string) {

	listen, err := net.Listen("tcp", port)
	if err != nil {
	}
	for {
		code := server.Accept(listen, codec.JsonType)

		h, b, _ := server.Read(*code)

		server.Handle(b)

		server.Send(*code, *h, *b)
	}
}

func (server *Server) Register(newStruct interface{}) bool {
	servicex := NewService(newStruct)

	server.mu.Lock()
	server.Services[servicex.Name] = servicex
	server.mu.Unlock()

	return true
}

func (server *Server) Accept(lis net.Listener, encodingType codec.Type) *codec.Codec {
	conn, err := lis.Accept()
	if err != nil {
		log.Println("server accept error:", err)
	}
	codefun := codec.NewCodecFuncMap[encodingType]
	request := codefun(conn)
	return &request
}

func (server *Server) Read(request codec.Codec) (*Header, *Body, error) {
	var h Header
	var err error
	err = request.ReadHeader(h)
	if err != nil {
		log.Println("server rpc read header error", err)
		return nil, nil, err
	}

	server.mu.Lock()
	ok := h.MagicNumber != server.MagicNumber
	server.mu.Unlock()

	if ok {
		log.Println("server rpc read header magicnumber is incorrect")
		err = errors.New("magicnumber is incorrect")
		return &h, nil, err
	}

	var body Body
	err = request.ReadBody(body)
	if err != nil {
		log.Println("server rpc read header error", err)
		return &h, nil, err
	}
	return &h, &body, nil
}

func (server *Server) Handle(body *Body) error {
	dot := strings.LastIndex(body.ServiceMethod, ".")
	serviceName, methodName := body.ServiceMethod[:dot], body.ServiceMethod[dot+1:]

	servicex, err := server.Services[serviceName]

	if err == false {
		log.Println("server rpc no service")
		body.err = errors.New("server rpc no service")
		return body.err
	}

	method, err := servicex.Method[methodName]
	if err == false {
		log.Println("server rpc service:", servicex.Name, "has  no ", methodName)
		body.err = errors.New("server rpc service:" + servicex.Name + "has  no " + methodName)

		return body.err
	}
	servicex.Call(method, body.Args, body.Reply)
	body.err = nil

	return nil
}

func (server *Server) Send(code codec.Codec, h Header, body Body) error {
	server.sending.Lock()
	err := code.Write(&h, &body)
	server.sending.Unlock()
	return err
}
