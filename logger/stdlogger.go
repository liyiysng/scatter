package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func init() {
	SetLogger(newLogger())
}

const depthStd = 2

const (
	// infoLog indicates Info severity.
	infoLog int = iota
	// warningLog indicates Warning severity.
	warningLog
	// errorLog indicates Error severity.
	errorLog
	// fatalLog indicates Fatal severity.
	fatalLog
)

// severityName contains the string representation of each severity.
var severityName = []string{
	infoLog:    "INFO",
	warningLog: "WARNING",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
}

// loggerT is the default logger used by logger.
type loggerT struct {
	m []*log.Logger
	v int
}

// NewLogger creates a Logger with the provided writers.
// Fatal logs will be written to errorW, warningW, infoW, followed by exit(1).
// Error logs will be written to errorW, warningW and infoW.
// Warning logs will be written to warningW and infoW.
// Info logs will be written to infoW.
func NewLogger(infoW, warningW, errorW io.Writer) Logger {
	return NewLoggerWithVerbosity(infoW, warningW, errorW, 0)
}

// NewLoggerWithVerbosity creates a loggerV2 with the provided writers and
// verbosity level.
func NewLoggerWithVerbosity(infoW, warningW, errorW io.Writer, v int) Logger {
	var m []*log.Logger
	m = append(m, log.New(infoW, severityName[infoLog]+": ", log.LstdFlags|log.Lshortfile))
	m = append(m, log.New(io.MultiWriter(infoW, warningW), severityName[warningLog]+": ", log.LstdFlags|log.Lshortfile))
	ew := io.MultiWriter(infoW, warningW, errorW) // ew will be used for error and fatal.
	m = append(m, log.New(ew, severityName[errorLog]+": ", log.LstdFlags|log.Lshortfile))
	m = append(m, log.New(ew, severityName[fatalLog]+": ", log.LstdFlags|log.Lshortfile))
	return &loggerT{m: m, v: v}
}

// newLogger creates a loggerV2 to be used as default logger.
// All logs are written to stderr.
func newLogger() Logger {
	errorW := ioutil.Discard
	warningW := ioutil.Discard
	infoW := ioutil.Discard

	logLevel := os.Getenv("SCATTER_GO_LOG_SEVERITY_LEVEL")

	logLevel = "INFO"

	switch logLevel {
	case "", "ERROR", "error": // If env is unset, set level to ERROR.
		errorW = os.Stderr
	case "WARNING", "warning":
		warningW = os.Stderr
	case "INFO", "info":
		infoW = os.Stderr
	}

	var v int
	vLevel := os.Getenv("SCATTER_GO_LOG_VERBOSITY_LEVEL")
	if vl, err := strconv.Atoi(vLevel); err == nil {
		v = vl
	} else {
		v = VALL
	}
	return NewLoggerWithVerbosity(infoW, warningW, errorW, v)
}

func (g *loggerT) Info(args ...interface{}) {
	g.m[infoLog].Output(depthStd, fmt.Sprint(args...))
}

func (g *loggerT) Infoln(args ...interface{}) {
	g.m[infoLog].Output(depthStd, fmt.Sprintln(args...))
}

func (g *loggerT) Infof(format string, args ...interface{}) {
	g.m[infoLog].Output(depthStd, fmt.Sprintf(format, args...))
}

func (g *loggerT) Warning(args ...interface{}) {
	g.m[warningLog].Output(depthStd, fmt.Sprint(args...))
}

func (g *loggerT) Warningln(args ...interface{}) {
	g.m[warningLog].Output(depthStd, fmt.Sprintln(args...))
}

func (g *loggerT) Warningf(format string, args ...interface{}) {
	g.m[warningLog].Output(depthStd, fmt.Sprintf(format, args...))
}

func (g *loggerT) Error(args ...interface{}) {
	g.m[errorLog].Output(depthStd, fmt.Sprint(args...))
}

func (g *loggerT) Errorln(args ...interface{}) {
	g.m[errorLog].Output(depthStd, fmt.Sprintln(args...))
}

func (g *loggerT) Errorf(format string, args ...interface{}) {
	g.m[errorLog].Output(depthStd, fmt.Sprintf(format, args...))
}

func (g *loggerT) Fatal(args ...interface{}) {
	g.m[fatalLog].Output(depthStd, fmt.Sprint(args...))
	os.Exit(1)
}

func (g *loggerT) Fatalln(args ...interface{}) {
	g.m[fatalLog].Output(depthStd, fmt.Sprintln(args...))
	os.Exit(1)
}

func (g *loggerT) Fatalf(format string, args ...interface{}) {
	g.m[fatalLog].Output(depthStd, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (g *loggerT) InfoDepth(depth int, args ...interface{}) {
	g.m[infoLog].Output(depth+1, fmt.Sprint(args...))
}

func (g *loggerT) WarningDepth(depth int, args ...interface{}) {
	g.m[warningLog].Output(depth+1, fmt.Sprint(args...))
}

func (g *loggerT) ErrorDepth(depth int, args ...interface{}) {
	g.m[errorLog].Output(depth+1, fmt.Sprint(args...))
}

func (g *loggerT) FatalDepth(depth int, args ...interface{}) {
	g.m[fatalLog].Output(depth+1, fmt.Sprint(args...))
	os.Exit(1)
}

func (g *loggerT) V(l int) bool {
	return l <= g.v
}
