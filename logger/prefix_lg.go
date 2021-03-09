package logger

import (
	"fmt"
)

const prefixLogDepth = 3

// prefixLogger Logging method on a nil logs without any prefix.
type prefixLogger struct {
	logger DepthLogger
	prefix string
}

func (g *prefixLogger) Info(args ...interface{}) {
	g.InfoDepth(prefixLogDepth, fmt.Sprint(args...))
}

func (g *prefixLogger) Infoln(args ...interface{}) {
	g.InfoDepth(prefixLogDepth, fmt.Sprintln(args...))
}

func (g *prefixLogger) Infof(format string, args ...interface{}) {
	g.InfoDepth(prefixLogDepth, fmt.Sprintf(format, args...))
}

func (g *prefixLogger) Warning(args ...interface{}) {
	g.WarningDepth(prefixLogDepth, fmt.Sprint(args...))
}

func (g *prefixLogger) Warningln(args ...interface{}) {
	g.WarningDepth(prefixLogDepth, fmt.Sprintln(args...))
}

func (g *prefixLogger) Warningf(format string, args ...interface{}) {
	g.WarningDepth(prefixLogDepth, fmt.Sprintf(format, args...))
}

func (g *prefixLogger) Error(args ...interface{}) {
	g.ErrorDepth(prefixLogDepth, fmt.Sprint(args...))
}

func (g *prefixLogger) Errorln(args ...interface{}) {
	g.ErrorDepth(prefixLogDepth, fmt.Sprintln(args...))
}

func (g *prefixLogger) Errorf(format string, args ...interface{}) {
	g.ErrorDepth(prefixLogDepth, fmt.Sprintf(format, args...))
}

func (g *prefixLogger) Fatal(args ...interface{}) {
	g.FatalDepth(prefixLogDepth, fmt.Sprint(args...))
}

func (g *prefixLogger) Fatalln(args ...interface{}) {
	g.FatalDepth(prefixLogDepth, fmt.Sprintln(args...))
}

func (g *prefixLogger) Fatalf(format string, args ...interface{}) {
	g.FatalDepth(prefixLogDepth, fmt.Sprintf(format, args...))
}

func (g *prefixLogger) InfoDepth(depth int, args ...interface{}) {
	args = append([]interface{}{g.prefix}, args)
	g.logger.InfoDepth(depth+1, fmt.Sprint(args...))
}

func (g *prefixLogger) WarningDepth(depth int, args ...interface{}) {
	args = append([]interface{}{g.prefix}, args)
	g.logger.WarningDepth(depth+1, fmt.Sprint(args...))
}

func (g *prefixLogger) ErrorDepth(depth int, args ...interface{}) {
	args = append([]interface{}{g.prefix}, args)
	g.logger.ErrorDepth(depth+1, fmt.Sprint(args...))
}

func (g *prefixLogger) FatalDepth(depth int, args ...interface{}) {
	args = append([]interface{}{g.prefix}, args)
	g.logger.Fatal(depth+1, fmt.Sprint(args...))
}

func (g *prefixLogger) V(l int) bool {
	return g.logger.V(l)
}

// NewPrefixLogger creates a prefix logger with the given prefix.
func NewPrefixLogger(logger DepthLogger, prefix string) DepthLogger {
	return &prefixLogger{logger: logger, prefix: prefix}
}
