package util

import (
	"sync/atomic"
	"time"
)

// AtomicBool bool 原子
// 默认为falses
type AtomicBool int32

// IsSet 是否被设置
func (b *AtomicBool) IsSet() bool { return atomic.LoadInt32((*int32)(b)) != 0 }

// SetTrue 设置为true
func (b *AtomicBool) SetTrue() { atomic.StoreInt32((*int32)(b), 1) }

// An AtomicDuration is an implementation of atomic duration.
type AtomicDuration int64

// NewAtomicDuration returns an AtomicDuration.
func NewAtomicDuration() *AtomicDuration {
	return new(AtomicDuration)
}

// ForAtomicDuration returns an AtomicDuration with given value.
func ForAtomicDuration(val time.Duration) *AtomicDuration {
	d := NewAtomicDuration()
	d.Set(val)
	return d
}

// CompareAndSwap compares current value with old, if equals, set the value to val.
func (d *AtomicDuration) CompareAndSwap(old, val time.Duration) bool {
	return atomic.CompareAndSwapInt64((*int64)(d), int64(old), int64(val))
}

// Load loads the current duration.
func (d *AtomicDuration) Load() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(d)))
}

// Set sets the value to val.
func (d *AtomicDuration) Set(val time.Duration) {
	atomic.StoreInt64((*int64)(d), int64(val))
}
