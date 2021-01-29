// Package logger 适配现有的日志库, loki , glog...
package logger

import (
	"os"
)

// GLogger is the logger used for the non-depth log functions.
var GLogger Logger

// GDepthLogger is the logger used for the depth log functions.
var GDepthLogger DepthLogger

// InfoDepth logs to the INFO log at the specified depth.
func InfoDepth(depth int, args ...interface{}) {
	if GDepthLogger != nil {
		GDepthLogger.InfoDepth(depth, args...)
	} else {
		GLogger.Infoln(args...)
	}
}

// WarningDepth logs to the WARNING log at the specified depth.
func WarningDepth(depth int, args ...interface{}) {
	if GDepthLogger != nil {
		GDepthLogger.WarningDepth(depth, args...)
	} else {
		GLogger.Warningln(args...)
	}
}

// ErrorDepth logs to the ERROR log at the specified depth.
func ErrorDepth(depth int, args ...interface{}) {
	if GDepthLogger != nil {
		GDepthLogger.ErrorDepth(depth, args...)
	} else {
		GLogger.Errorln(args...)
	}
}

// FatalDepth logs to the FATAL log at the specified depth.
func FatalDepth(depth int, args ...interface{}) {
	if GDepthLogger != nil {
		GDepthLogger.FatalDepth(depth, args...)
	} else {
		GLogger.Fatalln(args...)
	}
	os.Exit(1)
}

// Logger 不含调用堆栈深度
type Logger interface {
	// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
	Info(args ...interface{})
	// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
	Infoln(args ...interface{})
	// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
	Infof(format string, args ...interface{})
	// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
	Warning(args ...interface{})
	// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
	Warningln(args ...interface{})
	// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
	Warningf(format string, args ...interface{})
	// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	Error(args ...interface{})
	// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
	Errorln(args ...interface{})
	// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	Errorf(format string, args ...interface{})
	// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	// gRPC ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatal(args ...interface{})
	// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
	// gRPC ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatalln(args ...interface{})
	// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	// gRPC ensures that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatalf(format string, args ...interface{})
	// V reports whether verbosity level l is at least the requested verbose level.
	V(l int) bool
}

// DepthLogger 包含堆栈深度
type DepthLogger interface {
	Logger
	// InfoDepth logs to INFO log at the specified depth. Arguments are handled in the manner of fmt.Print.
	InfoDepth(depth int, args ...interface{})
	// WarningDepth logs to WARNING log at the specified depth. Arguments are handled in the manner of fmt.Print.
	WarningDepth(depth int, args ...interface{})
	// ErrorDetph logs to ERROR log at the specified depth. Arguments are handled in the manner of fmt.Print.
	ErrorDepth(depth int, args ...interface{})
	// FatalDepth logs to FATAL log at the specified depth. Arguments are handled in the manner of fmt.Print.
	FatalDepth(depth int, args ...interface{})
}

// V reports whether verbosity level l is at least the requested verbose level.
func V(l int) bool {
	return GLogger.V(l)
}

// SetLogger 设置全局日志对象
func SetLogger(l Logger) {
	if _, ok := l.(*componentData); ok {
		panic("cannot use component logger as gloable logger")
	}
	GLogger = l
	GDepthLogger, _ = l.(DepthLogger)
}
