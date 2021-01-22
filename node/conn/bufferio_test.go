package conn

import (
	"bufio"
	"io"
	"testing"
)

type buf20Bytes struct {
	r int
}

func (b *buf20Bytes) Read(p []byte) (n int, err error) {
	if len(p) < 20 {
		return 0, io.ErrShortBuffer
	}

	return 20, io.EOF

}

func TestBufferIOEOF(t *testing.T) {
	buf := &buf20Bytes{}
	bufIORead := bufio.NewReader(buf)
	_, err := bufIORead.Peek(20)
	if err != nil {
		t.Fatal(err)
		return
	}
	// buffio.Peek : 当给定N等于返回len(p)时,不会返回EOF,只有在下次读时返回EOF
	_, err = bufIORead.Peek(20)
	if err != nil {
		t.Fatal(err)
		return
	}
	_, err = bufIORead.Discard(20)
	if err != nil {
		t.Fatal(err)
		return
	}
	p := make([]byte, 1)
	_, err = bufIORead.Read(p)
	if err != nil {
		t.Fatal(err)
		return
	}
}
