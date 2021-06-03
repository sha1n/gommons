package test

import (
	"bufio"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmulatedStdinReaderRead(t *testing.T) {
	values := generateRandomValues()
	userInput := strings.Join(values, "\n")

	stdinReader := NewEmulatedStdinReader(userInput)
	reader := bufio.NewReader(stdinReader)

	assertReadValue := func(expected string) {
		var value string
		var err error
		value, err = reader.ReadString('\n')

		assert.NoError(t, err)
		assert.Equal(t, expected, strings.TrimSpace(value))
	}

	for _, v := range values {
		assertReadValue(v)
	}

	_, err := reader.ReadString('\n')

	assert.Error(t, err, "an error is expected after the expected content is fully consumed")
}

func generateRandomValues() []string {
	return []string{
		fmt.Sprint(rand.Intn(1000)),
		fmt.Sprint(rand.Intn(1000)),
		fmt.Sprint(rand.Intn(1000)),
	}
}
