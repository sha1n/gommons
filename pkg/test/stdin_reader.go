package test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var ErrUnexpectedRead = errors.New("unexpected stdin read")

type EmulatedStdinReader struct {
	buf           *bytes.Buffer
	lines         []string
	nextLineIndex int
}

func NewEmulatedStdinReader(content string) *EmulatedStdinReader {
	emulatedReader := &EmulatedStdinReader{
		buf:           new(bytes.Buffer),
		lines:         strings.Split(content, "\n"),
		nextLineIndex: 0,
	}

	return emulatedReader
}

func (s *EmulatedStdinReader) Read(buf []byte) (read int, err error) {
	if s.nextLineIndex < len(s.lines) {
		if s.buf.Len() == 0 {
			s.buf.Reset()
			fmt.Fprintf(s.buf, "%s\n", s.lines[s.nextLineIndex])
			s.nextLineIndex++
		}

		return s.buf.Read(buf)
	} else {
		return 0, ErrUnexpectedRead
	}
}
