package govtil

import (
	"bytes"
	"io"
	"sync"
)

// A closeable multi-threaded buffer
type CloseableBuffer struct {
	buf *bytes.Buffer	// don't use embedding, since I don't want to refine
						// every Read*() and Write*() method, just the basic ones
	closed bool
	sync.Cond
}

func NewCloseableBuffer() *CloseableBuffer {
	return &CloseableBuffer{&bytes.Buffer{}, false, *sync.NewCond(&sync.Mutex{})}
}

// After Close() is called, Read() and Write() will fail with err == io.EOF
func (cb *CloseableBuffer) Close() {
	cb.L.Lock()
	defer cb.L.Unlock()
	cb.closed = true
}

func (cb *CloseableBuffer) Closed() bool {
	cb.L.Lock()
	defer cb.L.Unlock()
	return cb.closed
}

// Refined Read() does exactly one of the following:
// 	1) gets data if there is any in the buffer
//	2) waits for data if there isn't any, then tries again
//	3) returns (0, io.EOF) if buffer is closed
func (cb *CloseableBuffer) Read(data []byte) (n int, err error) {
	cb.L.Lock()
	defer cb.L.Unlock()
	for n, err = cb.buf.Read(data); n == 0 ; {
		if cb.closed {
			return 0, io.EOF
		}
		cb.Wait()
	}
	return
}

// Refined Write() writes to buffer if open and signals waiting reads
func (cb *CloseableBuffer) Write(data []byte) (n int, err error) {
	cb.L.Lock()
	defer cb.L.Unlock()
	if cb.closed {
		return 0, io.EOF
	}
	n, err = cb.buf.Write(data)
	cb.Signal()
	return
}

func (cb *CloseableBuffer) Len() int {
	cb.L.Lock()
	defer cb.L.Unlock()
	return cb.buf.Len()
}
