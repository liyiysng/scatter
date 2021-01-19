package message

import (
	"encoding/binary"
	"io"
)

// Packet 一个包含完整消息的数据包
type Packet struct {
	msgType MsgType
	length  int32
	data    []byte
}

// ReadService 读取服务名
// 由于各个平台的字符串表达方式不同,采用 1byte(str length)+str
func ReadService(r io.ByteScanner) (srv string, err error) {
	var strLen uint8
	err = binary.Read(r, binary.BigEndian, &strLen)
	if err != nil {
		return "", err
	}

}
