package io

import (
	"bytes"
	"io"
	"sync"
)

// ProbedWriter wraps an io.Writer and maintains the last N bytes written to it.
type ProbedWriter struct {
	w       io.Writer
	buf     *bytes.Buffer
	maxSize int
	mu      sync.Mutex
}

// NewProbedWriter creates a new ProbedWriter.
func NewProbedWriter(w io.Writer, n int) *ProbedWriter {
	return &ProbedWriter{
		w:       w,
		buf:     bytes.NewBuffer(make([]byte, 0, n)),
		maxSize: n,
	}
}

// Write writes p to the underlying io.Writer and keeps the last N bytes in a buffer.
func (l *ProbedWriter) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	n, err = l.w.Write(p)
	if err != nil {
		return n, err
	}

	// Write to our internal buffer.
	_, _ = l.buf.Write(p[:n])

	// If our buffer exceeds maxSize, trim it.
	if l.buf.Len() > l.maxSize {
		overflow := l.buf.Len() - l.maxSize
		_ = l.buf.Next(overflow)
	}

	return n, nil
}

// LastN returns the last N bytes written.
func (l *ProbedWriter) Bytes() []byte {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.Bytes()
}

// String returns the last N bytes written as a string.
func (l *ProbedWriter) String() string {
	return string(l.Bytes())
}

func (l *ProbedWriter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buf.Reset()
}
