package util

import (
	"bytes"
	"sync"
)

var bp sync.Pool

func init() {
	bp.New = func() interface{} {
		return &bytes.Buffer{}
	}
}

// BufferPoolGet 获得一个buffer
func BufferPoolGet() *bytes.Buffer {
	return bp.Get().(*bytes.Buffer)
}

// BufferPoolPut 回收一个buffer
func BufferPoolPut(b *bytes.Buffer) {
	b.Reset()
	bp.Put(b)
}
