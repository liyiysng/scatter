package util

import "sync/atomic"

// AtomicBool bool 原子
// 默认为falses
type AtomicBool int32

// IsSet 是否被设置
func (b *AtomicBool) IsSet() bool { return atomic.LoadInt32((*int32)(b)) != 0 }

// SetTrue 设置为true
func (b *AtomicBool) SetTrue() { atomic.StoreInt32((*int32)(b), 1) }
