package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/liyiysng/scatter/constants"
	"github.com/liyiysng/scatter/encoding"
	"github.com/liyiysng/scatter/util"
)

type packageWriter interface {
	io.ByteWriter
	io.Writer
}

type packageReader interface {
	io.ByteReader
	io.Reader
}

var pkgPool sync.Pool

func init() {
	pkgPool.New = func() interface{} {
		return &Packet{}
	}
}

// PackagePoolGet 获得一个 packet
func PackagePoolGet() *Packet {
	return pkgPool.Get().(*Packet)
}

// PackagePoolPut 回收一个 packet
func PackagePoolPut(p *Packet) {
	p.Reset()
	pkgPool.Put(p)
}

// Packet 一个包含完整消息的数据包
type Packet struct {
	MsgType MsgType
	MsgOpt  MsgOpt
	Data    []byte
}

// ReadFrom 从reader中读取package
func (p *Packet) ReadFrom(r packageReader, compressor encoding.Compressor, maxLength int) (n int, err error) {

	bType, err := r.ReadByte()
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return n, err
	}
	n++

	p.MsgType = MsgType(bType)

	bOpt, err := r.ReadByte()
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return n, err
	}
	n++

	p.MsgOpt = MsgOpt(bOpt)

	dataLength := int32(-1)
	err = binary.Read(r, binary.BigEndian, &dataLength)
	if err != nil {
		return n, err
	}
	if dataLength > int32(maxLength) || dataLength < 0 {
		return n, constants.ErrMsgTooLager
	}
	n += 4

	if dataLength > 0 {
		buf := make([]byte, int(dataLength))
		read, err := io.ReadFull(r, buf)
		n += read
		if err != nil {
			return n, err
		}
		if p.MsgOpt&COMPRESS > 0 { // 解压
			r, err := compressor.Decompress(bytes.NewBuffer(buf))
			if err != nil {
				return n, err
			}
			buf, err = ioutil.ReadAll(r)
			if err != nil {
				return n, err
			}
		}
		p.Data = buf
	} else {
		p.Data = nil
	}
	return
}

// WriteTo 写入一个package
func (p *Packet) WriteTo(w packageWriter, compresser encoding.Compressor, maxLength int) (n int, err error) {

	length := len(p.Data)
	if length > maxLength {
		return 0, constants.ErrMsgTooLager
	}

	bufToWrite := p.Data

	if length > 0 && p.MsgOpt&COMPRESS > 0 { // 压缩
		compressBuf := util.BufferPoolGet()
		defer util.BufferPoolPut(compressBuf)
		wr, err := compresser.Compress(compressBuf)
		if err != nil {
			return 0, err
		}
		_, err = wr.Write(p.Data)
		if err != nil {
			return 0, err
		}
		bufToWrite = compressBuf.Bytes()
	}

	lenghtToWrite := int32(len(bufToWrite))

	err = w.WriteByte(byte(p.MsgType))
	if err != nil {
		return
	}
	n++
	err = w.WriteByte(byte(p.MsgOpt))
	if err != nil {
		return
	}
	n++
	err = binary.Write(w, binary.BigEndian, &lenghtToWrite)
	if err != nil {
		return
	}
	n += 4
	if lenghtToWrite > 0 {
		written, err := w.Write(bufToWrite)
		if err != nil {
			return n, err
		}
		n += written
	}
	return
}

// Reset 重置Packet
func (p *Packet) Reset() {
	p.Data = nil
}

// ReadService 读取服务名
// 由于各个平台的字符串表达方式不同,采用 1byte(str length)+str
func ReadService(r io.Reader) (service string, err error) {
	var strLen uint8
	err = binary.Read(r, binary.BigEndian, &strLen)
	if err != nil {
		return "", err
	}
	if strLen > constants.MaxMessageServiceName {
		return "", fmt.Errorf("[ReadService] invalid service length %d", strLen)
	}
	buf := make([]byte, strLen)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
