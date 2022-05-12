package codec

import (
	// "fmt"
	"io"
)

//将这个
type Type string

type NewCodecFunc func(io.ReadWriteCloser) Codec //函数指针（模板）
const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // not implemented
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	//分配内存
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[JsonType] = NewJsonCodec
	NewCodecFuncMap[GobType] = NewGobCodec
}

type Header struct {
	ServiceMethod string // format "Service.Method"
	Seq           uint64 // sequence number chosen by client
	Error         string
}

//Can use different impl
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}
