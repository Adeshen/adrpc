package adrpc

import (
	"adrpc/codec"
	"fmt"
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
		fmt.Print("server listen create failed")
	}
	for {
		code := server.Accept(listen, codec.JsonType)

		h, b, s, m := server.Read(*code)
		fmt.Println(b.Args)
		if h == nil {
			log.Println("head is nil   ")
			continue
		}

		if b == nil {
			log.Println("from " + string(h.Clientid) + " body is nil")
			continue
		}
		server.GetReply(s, m, b)

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

func (server *Server) Read(request codec.Codec) (*Header, *Body, *service, *methodType) {
	var h Header
	var err error
	err = request.ReadHeader(&h)
	if err != nil {
		log.Println("server rpc read header error", err)
		return nil, nil, nil, nil
	}

	server.mu.Lock()
	ok := h.MagicNumber != server.MagicNumber
	server.mu.Unlock()

	if ok {
		log.Println("server rpc read header magicnumber is incorrect")
		// err = errors.New("magicnumber is incorrect")

		return &h, nil, nil, nil
	}

	var body Body
	service, method := server.Find(h.ServiceMethod)
	server.PreBody(method, &body)
	body.Args = method.newArgv()
	argvi := body.Args.Interface()

	// if body.Args.Type().Kind() != reflect.Ptr {
	// 	argvi = body.Args.Addr().Interface()
	// }
	err = request.ReadBody(&argvi)
	// err = request.ReadBody(&body)

	if err != nil {
		log.Println("server rpc read header error", err)
		return &h, &body, service, method
	}

	return &h, &body, service, method
}

func (server *Server) Find(serviceMethod string) (*service, *methodType) {
	dot := strings.LastIndex(serviceMethod, ".")
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	servicex, err := server.Services[serviceName]
	if err == false {
		return nil, nil
	}
	method, err := servicex.Method[methodName]
	if err == false {
		return servicex, nil
	}
	return servicex, method
}

func (server *Server) PreBody(method *methodType, body *Body) *Body {
	body.Args = method.newArgv()
	body.Reply = method.newReplyv()
	return body
}

func (server *Server) GetReply(servicex *service, method *methodType, prebody *Body) {
	servicex.Call(method, prebody.Args, prebody.Reply)
	prebody.err = nil
}

// func (server *Server) Handle(body *Body) error {
// 	fmt.Println(*body)
// 	dot := strings.LastIndex(body.ServiceMethod, ".")
// 	serviceName, methodName := body.ServiceMethod[:dot], body.ServiceMethod[dot+1:]

// 	servicex, err := server.Services[serviceName]

// 	if err == false {
// 		log.Println("server rpc no service")
// 		body.err = errors.New("server rpc no service")
// 		return body.err
// 	}

// 	method, err := servicex.Method[methodName]
// 	if err == false {
// 		log.Println("server rpc service:", servicex.Name, "has  no ", methodName)
// 		body.err = errors.New("server rpc service:" + servicex.Name + "has  no " + methodName)

// 		return body.err
// 	}

// 	return nil
// }

func (server *Server) Send(code codec.Codec, h Header, body Body) error {
	server.sending.Lock()
	replyiv := body.Reply.Interface()
	err := code.Write(&h, &replyiv)
	server.sending.Unlock()
	return err
}
