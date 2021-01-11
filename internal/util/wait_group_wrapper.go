package util

import (
	"sync"
)

// WaitGroupWrapper 检测goroutine
type WaitGroupWrapper struct {
	sync.WaitGroup
}

// Wrap 封装go语句,检测goroutine退出
func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
