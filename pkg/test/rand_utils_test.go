package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomStrings(t *testing.T) {
	assert.Eventually(
		t,
		func() bool { return len(RandomStrings()) != len(RandomStrings()) },
		time.Second,
		time.Millisecond,
	)
}

func TestRandomString(t *testing.T) {
	assert.Eventually(
		t,
		func() bool { return RandomString() != RandomString() },
		time.Second,
		time.Millisecond,
	)
}

func TestRandomBool(t *testing.T) {
	assert.Eventually(
		t,
		func() bool { return RandomBool() != RandomBool() },
		time.Second,
		time.Nanosecond,
	)
}

func TestRandomUint(t *testing.T) {
	assert.Eventually(
		t,
		func() bool { return RandomUint() != RandomUint() },
		time.Second,
		time.Millisecond,
	)
}
