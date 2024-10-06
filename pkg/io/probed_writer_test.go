package io

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/sha1n/gommons/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestImplementsWriter(t *testing.T) {
	var writer interface{} = &ProbedWriter{}
	_, ok := writer.(io.Writer)

	assert.True(t, ok, "ProbedWriter does not implement io.Writer")
}

func TestEmpty(t *testing.T) {
	wp := NewUnlimitedProbedWriter(new(bytes.Buffer))

	assert.Equal(t, wp.Len(), 0)
	assert.Equal(t, len(wp.Bytes()), 0)
	assert.Equal(t, wp.String(), "")
}

func TestReset(t *testing.T) {
	wp := NewUnlimitedProbedWriter(new(bytes.Buffer))
	wp.Write([]byte(test.RandomString()))

	wp.Reset()

	assert.Equal(t, wp.Len(), 0)
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

		assert.Equal(t, wp.Len(), len(tc.expectedOutput))
		assert.Equal(t, len(wp.Bytes()), len(tc.expectedOutput))
		assert.Equal(t, wp.String(), tc.expectedOutput)

		wp.Reset()
	}
}

func TestWriteProbeBufferUnlimitedContent(t *testing.T) {
	wp := NewUnlimitedProbedWriter(new(bytes.Buffer))

	testCases := []struct {
		input string
	}{
		{""},
		{generateRandomString(1, 1024)},
		{generateRandomString(1024, 1024*1024*10)},
	}

	for _, tc := range testCases {
		_, err := wp.Write([]byte(tc.input))
		assert.NoError(t, err)

		assert.Equal(t, wp.Len(), len(tc.input))
		assert.Equal(t, len(wp.Bytes()), len(tc.input))
		assert.Equal(t, wp.String(), tc.input)

		wp.Reset()
	}
}

func TestWriteOutputWriter(t *testing.T) {
	outputWriter := new(bytes.Buffer)
	wp := NewProbedWriter(outputWriter, 1)

	wp.Write([]byte("Hello, world!"))

	assert.Equal(t, "Hello, world!", outputWriter.String())
}

// GenerateRandomString generates a random string of length between minLen and maxLen.
func generateRandomString(minLen, maxLen int) string {
	// FIXME shai: promote to test package
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Define the characters to be used in the random string
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Calculate the random length between minLen and maxLen
	length := random.Intn(maxLen-minLen+1) + minLen

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
