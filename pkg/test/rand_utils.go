package test

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

// RandomStrings returns a slice of random strings
func RandomStrings() []string {
	values := []string{}
	for i := 0; i < rand.Intn(10); i++ {
		values = append(values, RandomString())
	}

	return values
}

// RandomString returns a random UUID based string...
func RandomString() string {
	uid, _ := uuid.NewRandom()
	return fmt.Sprintf("str-%s", uid.String())
}

// RandomBool ...
func RandomBool() bool {
	return rand.Intn(2)%2 == 0
}

// RandomUint ...
func RandomUint() uint {
	return uint(rand.Uint32())
}
