package io

import (
	"bytes"
	"io"
	"testing"

	"github.com/sha1n/gommons/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestImplementsWriter(t *testing.T) {
	var writer interface{} = &ProbedWriter{}
	_, ok := writer.(io.Writer)

	assert.True(t, ok, "ProbedWriter does not implement io.Writer")
}

func TestEmpty(t *testing.T) {
	wp := NewProbedWriter(new(bytes.Buffer), 10)

	assert.Equal(t, len(wp.Bytes()), 0)
	assert.Equal(t, wp.String(), "")
}

func TestReset(t *testing.T) {
	wp := NewProbedWriter(new(bytes.Buffer), 10)
	wp.Write([]byte(test.RandomString()))

	wp.Reset()

	assert.Equal(t, len(wp.Bytes()), 0)
	assert.Equal(t, wp.String(), "")
}

func TestWriteProbeBufferContent(t *testing.T) {
	probeBufferSize := 10
	wp := NewProbedWriter(new(bytes.Buffer), probeBufferSize)

	testCases := []struct {
		input          string
		expectedOutput string
	}{
		{"", ""},
		{"Hell", "Hell"},
		{"Hello", "Hello"},
		{"Hello, world!", "lo, world!"},
		{"こんにちは、世界!", "、世界!"}, // 3x3 bytes + 1 byte for the exclamation mark
	}

	for _, tc := range testCases {
		_, err := wp.Write([]byte(tc.input))
		assert.NoError(t, err)

		assert.Equal(t, len(wp.Bytes()), len(tc.expectedOutput))
		assert.Equal(t, wp.String(), tc.expectedOutput)

		wp.Reset()
	}
}

func TestWriteOutputWriter(t *testing.T) {
	outputWriter := new(bytes.Buffer)
	wp := NewProbedWriter(outputWriter, 1)

	wp.Write([]byte("Hello, world!"))

	assert.Equal(t, "Hello, world!", outputWriter.String())
}
