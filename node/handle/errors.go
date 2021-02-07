package handle

import "fmt"

// ICustomError 用户定义错误
// 不会导致链接关闭(如:通知卡牌数量不足等...)
type ICustomError interface {
	error
	customErrorMark()
}

// CustomError implements error interface
type CustomError struct {
	err string
}

func (e *CustomError) Error() string {
	return e.err
}

func (e *CustomError) customErrorMark() {
	panic("customErrorMark")
}

// NewCustomErrorf 创建自定义错误
func NewCustomErrorf(format string, args ...interface{}) ICustomError {
	return NewCustomError(fmt.Sprintf(format, args...))
}

// NewCustomError 创建自定义错误
func NewCustomError(err string) ICustomError {
	return &CustomError{
		err: err,
	}
}

// NewCustomErrorWithError 创建自定义错误
func NewCustomErrorWithError(err error) *CustomError {
	return &CustomError{
		err: err.Error(),
	}
}

// ICriticalError 关键错误
// 导致session关闭
type ICriticalError interface {
	error
	criticalErrorMark()
}

// CriticalError implements error interface
type CriticalError struct {
	err string
}

func (e *CriticalError) Error() string {
	return e.err
}

func (e *CriticalError) criticalErrorMark() {
	panic("criticalErrorMark")
}

// NewCriticalErrorf 创建关键性错误
func NewCriticalErrorf(format string, args ...interface{}) ICriticalError {
	return NewCriticalError(fmt.Sprintf(format, args...))
}

// NewCriticalError 创建关键性错误
func NewCriticalError(err string) ICriticalError {
	return &CriticalError{
		err: err,
	}
}

// NewCriticalErrorWithError 创建关键性错误
func NewCriticalErrorWithError(err error) ICriticalError {
	return &CriticalError{
		err: err.Error(),
	}
}
