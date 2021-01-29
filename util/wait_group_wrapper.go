package util

import (
	"runtime"
	"sync"
)

// WaitGroupWrapper 检测goroutine
type WaitGroupWrapper struct {
	sync.WaitGroup
}

// Wrap 封装go语句,检测goroutine退出
func (w *WaitGroupWrapper) Wrap(cb func(), errOutput func(format string, v ...interface{})) {
	w.Add(1)
	go func() {
		if errOutput != nil {
			defer func() {
				if err := recover(); err != nil {
					const size = 64 << 10
					buf := make([]byte, size)
					buf = buf[:runtime.Stack(buf, false)]
					errOutput("goroutine panic %v\n%s", err, buf)
				}
			}()
		}

		defer func() {
			w.Done()
		}()
		cb()
	}()
}
