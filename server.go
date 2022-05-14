package adrpc

import (
	"adrpc/codec"
	"fmt"
	"log"
	"net"
	"reflect"
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
		req := server.NewRequest(h, b, s, m)

		req.Args = req.Methodt.newArgv()
		req.Reply = req.Methodt.newReplyv()

		// make sure that argvi is a pointer, ReadBody need a pointer as parameter
		//原来真的是可以先预设反射类型，然后把用编码器把interface编码成字节数组，然后解码到反射的输入接口
		argvi := req.Args.Interface()
		if req.Args.Type().Kind() != reflect.Ptr {
			argvi = req.Args.Addr().Interface()
		}
		if err = (*code).ReadBody(argvi); err != nil {
			log.Println("rpc server: read body err:", err)
			return
		}
		// fmt.Println("argvi:", argvi)
		// fmt.Println("接受到的args:", req.Args.Interface())
		server.GetReply(req)
		server.Send(*code, req)
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
	service, method := server.Find(h.ServiceMethod)
	body := Body{}
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

func (server *Server) NewRequest(h *Header, b *Body, s *service, m *methodType) *Request {
	req := Request{
		RH:      h,
		Args:    m.newArgv(), //而此时有参数类型是interface{},我想把它转化想要的已经存储在m.ArgType的类型
		Reply:   m.newReplyv(),
		Service: s,
		Methodt: m,
	}

	// fmt.Println("反射args：", req.Args.Interface())
	return &req
}

func (server *Server) GetReply(request *Request) {
	// servicex.Call(method, request.Args, request.Reply)
	// prebody.err = nil
	request.Service.Call(request.Methodt, request.Args, request.Reply)
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

func (server *Server) Send(code codec.Codec, req *Request) error {
	b := Body{
		Args:  req.Args.Interface(),
		Reply: req.Reply.Interface(),
		err:   nil,
	}
	server.sending.Lock()
	err := code.Write(&req.RH, &b)
	server.sending.Unlock()
	return err
}
