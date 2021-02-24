package textlog

import (
	"fmt"
	"io/ioutil"

	"github.com/liyiysng/scatter/logger"
)

var (
	myLog = logger.Component("textlog")
)

// NewTempFileSink creates a temp file and returns a Sink that writes to this
// file.
func NewTempFileSink() (Sink, error) {
	// Two other options to replace this function:
	// 1. take filename as input.
	// 2. export NewBufferedSink().
	tempFile, err := ioutil.TempFile("./tmp", "scatter_textlog_*.log")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	return NewBufferedSink(tempFile), nil
}
