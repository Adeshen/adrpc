package adrpc

import (
	"adrpc/codec"
	"errors"
	"log"
	"net"
	"time"
)

func RomoteCall(serviceMethod string, arg interface{}, reply interface{}) {

}

func NewClient(addr string, CodecType codec.Type, Conntime time.Duration) (*Client, error) {
	conn, err := net.DialTimeout("tcp", addr, Conntime)
	if err != nil {
		log.Fatal("timeout")
		return nil, errors.New("client create defeat")
	}
	codefun := codec.NewCodecFuncMap[CodecType]
	client := Client{
		NetIO:    codefun(conn),
		Seq:      1,
		Callchan: make(chan *Call),
		Pending:  map[int]*Call{},
	}
	return &client, nil
}

//
func (client *Client) AddCall(serviceMethod string, args interface{}, reply interface{}) *Call {
	var call Call
	call = Call{
		Method: serviceMethod,
		Args:   args,
		Reply:  reply,
	}

	client.Mu.Lock()
	call.Seq = client.Seq
	client.Seq++
	client.Pending[call.Seq] = &call
	client.Mu.Unlock()
	return &call
}

type B struct {
	Args  interface{}
	Reply interface{}
	err   error
}

func (client *Client) Send(id uint64, MagicNumber uint64, call *Call) {
	h := Header{
		MagicNumber:   MagicNumber,
		Seq:           call.Seq,
		Clientid:      id,
		ServiceMethod: call.Method,
	}

	// b := Body{
	// 	Args:  reflect.ValueOf(call.Args),
	// 	Reply: reflect.ValueOf(call.Reply),
	// 	err:   nil,
	// }
	b := B{
		Args: (call.Args),
	}

	client.NetIO.Write(&h, &b)
}

func (client *Client) Receive() (*Header, *Body) {
	var h Header
	client.NetIO.ReadBody(&h)
	var b Body
	client.NetIO.ReadBody(&b)
	return &h, &b
}

func (client *Client) Close() {
	client.NetIO.Close()
}
