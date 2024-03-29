package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"sync"

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

// PacketOpt 消息选项
type PacketOpt byte

const (
	// DEFAULTPOPT 默认选项
	DEFAULTPOPT PacketOpt = 0x00
	// COMPRESS 是否压缩
	COMPRESS PacketOpt = 0x01
)

// IPacketOption 包选项
type IPacketOption interface {
	Apply(*PacketOpt)
}

type funcPacketOption struct {
	f func(*PacketOpt)
}

func (fdo *funcPacketOption) Apply(do *PacketOpt) {
	fdo.f(do)
}

func newFuncPacketOption(f func(*PacketOpt)) *funcPacketOption {
	return &funcPacketOption{
		f: f,
	}
}

// WithCompressOpt 压缩选项
func WithCompressOpt() IPacketOption {
	return newFuncPacketOption(func(o *PacketOpt) {
		*o |= COMPRESS
	})
}

// Packet 一个包含完整消息的数据包
type Packet struct {
	PacketOpt PacketOpt
	Data      []byte
}

// ReadFrom 从reader中读取package
func (p *Packet) ReadFrom(r packageReader, compressor encoding.Compressor, maxLength int) (n int, err error) {

	bOpt, err := r.ReadByte()
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return n, err
	}
	n++

	p.PacketOpt = PacketOpt(bOpt)

	dataLength := int32(-1)
	err = binary.Read(r, binary.BigEndian, &dataLength)
	if err != nil {
		return n, err
	}
	if dataLength > int32(maxLength) || dataLength < 0 {
		return n, ErrMsgTooLager
	}
	n += 4

	if dataLength > 0 {
		buf := make([]byte, int(dataLength))
		read, err := io.ReadFull(r, buf)
		n += read
		if err != nil {
			return n, err
		}
		if p.PacketOpt&COMPRESS > 0 { // 解压
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
		err = errors.New("read empty message")
	}
	return
}

// WriteTo 写入一个package
func (p *Packet) WriteTo(w packageWriter, compresser encoding.Compressor, maxLength int) (n int, err error) {

	length := len(p.Data)
	if length > maxLength {
		return 0, ErrMsgTooLager
	}

	bufToWrite := p.Data

	if length > 0 && p.PacketOpt&COMPRESS > 0 { // 压缩

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
		err = wr.Close()
		if err != nil {
			return 0, err
		}
		bufToWrite = compressBuf.Bytes()
	}

	lenghtToWrite := int32(len(bufToWrite))

	err = w.WriteByte(byte(p.PacketOpt))
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
