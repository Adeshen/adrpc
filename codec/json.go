package codec

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

type JsonCodec struct {
	conn io.ReadWriteCloser //存储连接
	buf  *bufio.Writer
	dec  *json.Decoder
	enc  *json.Encoder
}

func (c *JsonCodec) Close() error {
	return c.conn.Close()
}

//解码头部  docode header
func (c *JsonCodec) ReadHeader(h interface{}) error {
	return c.dec.Decode(h)
}

//解码body  decode body
func (c *JsonCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

//
func (c *JsonCodec) Write(h interface{}, body interface{}) (err error) {
	//最后执行  压入栈钟倒序执行
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	//编码失败
	if err := c.enc.Encode(h); err != nil {
		log.Println("rpc codec:T", h)
		log.Println("rpc codec: json error encoding header:", err)
		return err
	}
	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec: json error encoding body:", err)
		return err
	}
	return nil
}

func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &JsonCodec{
		conn: conn,
		buf:  buf,
		dec:  json.NewDecoder(conn), //这里可以直接编写进去
		enc:  json.NewEncoder(buf),
	}
}
