// Package lg short for "log"
package lg

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// LogLevel 日志显示等级
type LogLevel int

// Get 获取日志等级
func (l *LogLevel) Get() interface{} { return *l }

// Set 设置日志等级
func (l *LogLevel) Set(s string) error {
	lvl, err := ParseLogLevel(s)
	if err != nil {
		return err
	}
	*l = lvl
	return nil
}

// String
func (l *LogLevel) String() string {
	switch *l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return "invalid"
}

// ParseLogLevel 从给定字符串解析日志等级
func ParseLogLevel(levelstr string) (LogLevel, error) {
	switch strings.ToLower(levelstr) {
	case "debug":
		return DEBUG, nil
	case "info":
		return INFO, nil
	case "warn":
		return WARN, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	}
	return 0, fmt.Errorf("invalid log level '%s' (debug, info, warn, error, fatal)", levelstr)
}

const (
	// DEBUG 调试模式 , 显示所有信息
	DEBUG = LogLevel(1)
	// INFO 显示程序信息
	INFO = LogLevel(2)
	// WARN 警告信息
	WARN = LogLevel(3)
	// ERROR 错误信息
	ERROR = LogLevel(4)
	// FATAL 致命错误,可能导致宕机
	FATAL = LogLevel(5)
)

// AppLogFunc 上层应用日志输出回调类型
type AppLogFunc func(lvl LogLevel, f string, args ...interface{})

// Logger 日志接口
// 所有日志实体都必须实现该接口
type Logger interface {
	// 日志输出
	// calldepth : 堆栈深度
	// s : 日志内容
	Output(maxdepth int, s string) error
}

// NilLogger 空日志实体
type NilLogger struct{}

// Output 实现 Logger
func (l NilLogger) Output(maxdepth int, s string) error {
	return nil
}

// Logf 格式化输出
func Logf(logger Logger, cfgLevel LogLevel, msgLevel LogLevel, f string, args ...interface{}) {
	if cfgLevel > msgLevel {
		return
	}
	logger.Output(3, fmt.Sprintf(msgLevel.String()+": "+f, args...))
}

// LogFatal 致命错误输出
func LogFatal(prefix string, f string, args ...interface{}) {
	logger := log.New(os.Stderr, prefix, log.Ldate|log.Ltime|log.Lmicroseconds)
	Logf(logger, FATAL, FATAL, f, args...)
	os.Exit(1)
}
