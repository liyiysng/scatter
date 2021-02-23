package textlog

import (
	"bufio"
	"io"
	"sync"
	"time"

	"github.com/liyiysng/scatter/node/session"
)

var (
	// DefaultSink is the sink where the logs will be written to. It's exported
	DefaultSink Sink = &noopSink{}
)

// Sink writes log entry into the binary log sink.
//
// sink is a copy of the exported binarylog.Sink, to avoid circular dependency.
type Sink interface {
	// Write will be called to write the log entry into the sink.
	//
	// It should be thread-safe so it can be called in parallel.
	Write(*session.MsgInfo) error
	// Close will be called when the Sink is replaced by a new Sink.
	Close() error
}

type noopSink struct{}

func (ns *noopSink) Write(*session.MsgInfo) error { return nil }
func (ns *noopSink) Close() error                 { return nil }

// newWriterSink creates a binary log sink with the given writer.
//
// No buffer is done, Close() doesn't try to close the writer.
func newWriterSink(w io.Writer) Sink {
	return &writerSink{out: w}
}

type writerSink struct {
	out io.Writer
}

func (ws *writerSink) Write(info *session.MsgInfo) error {
	if _, err := ws.out.Write([]byte(info.String())); err != nil {
		return err
	}
	if _, err := ws.out.Write([]byte("\t")); err != nil {
		return err
	}
	return nil
}

func (ws *writerSink) Close() error { return nil }

type bufferedSink struct {
	mu     sync.Mutex
	closer io.Closer
	out    Sink          // out is built on buf.
	buf    *bufio.Writer // buf is kept for flush.

	writeStartOnce sync.Once
	writeTicker    *time.Ticker
}

func (fs *bufferedSink) Write(info *session.MsgInfo) error {
	// Start the write loop when Write is called.
	fs.writeStartOnce.Do(fs.startFlushGoroutine)
	fs.mu.Lock()
	if err := fs.out.Write(info); err != nil {
		fs.mu.Unlock()
		return err
	}
	fs.mu.Unlock()
	return nil
}

const (
	bufFlushDuration = 60 * time.Second
)

func (fs *bufferedSink) startFlushGoroutine() {
	fs.writeTicker = time.NewTicker(bufFlushDuration)
	go func() {
		for range fs.writeTicker.C {
			fs.mu.Lock()
			if err := fs.buf.Flush(); err != nil {
				myLog.Warningf("failed to flush to Sink: %v", err)
			}
			fs.mu.Unlock()
		}
	}()
}

func (fs *bufferedSink) Close() error {
	if fs.writeTicker != nil {
		fs.writeTicker.Stop()
	}
	fs.mu.Lock()
	if err := fs.buf.Flush(); err != nil {
		myLog.Warningf("failed to flush to Sink: %v", err)
	}
	if err := fs.closer.Close(); err != nil {
		myLog.Warningf("failed to close the underlying WriterCloser: %v", err)
	}
	if err := fs.out.Close(); err != nil {
		myLog.Warningf("failed to close the Sink: %v", err)
	}
	fs.mu.Unlock()
	return nil
}

// NewBufferedSink creates a binary log sink with the given WriteCloser.
//
// Write() marshals the proto message and writes it to the given writer. Each
// message is prefixed with a 4 byte big endian unsigned integer as the length.
//
// Content is kept in a buffer, and is flushed every 60 seconds.
//
// Close closes the WriteCloser.
func NewBufferedSink(o io.WriteCloser) Sink {
	bufW := bufio.NewWriter(o)
	return &bufferedSink{
		closer: o,
		out:    newWriterSink(bufW),
		buf:    bufW,
	}
}
