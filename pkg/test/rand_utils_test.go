package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomStrings(t *testing.T) {
	assert.Eventually(
		t,
		func() bool {
			return len(RandomStrings()) != len(RandomStrings()) //nolint:staticcheck
		},
		time.Second,
		time.Millisecond,
	)
}

func TestRandomString(t *testing.T) {
	assert.Eventually(
		t,
		func() bool {
			return RandomString() != RandomString() //nolint:staticcheck
		},
		time.Second,
		time.Millisecond,
	)
}

func TestRandomBool(t *testing.T) {
	assert.Eventually(
		t,
		func() bool {
			return RandomBool() != RandomBool() //nolint:staticcheck
		},
		time.Second,
		time.Nanosecond,
	)
}

func TestRandomUint(t *testing.T) {
	assert.Eventually(
		t,
		func() bool {
			return RandomUint() != RandomUint() //nolint:staticcheck
		},
		time.Second,
		time.Millisecond,
	)
}
