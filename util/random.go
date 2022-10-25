package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random int64 between min and max.
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder

	// Gets random chars from the alphabet
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(c) // appends to sb
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(8)
}

// RandomBalance generates a random amount of money
func RandomBalance() int64 {
	return RandomInt(0, 1000)
}

//RandomCurrency generates a random currency symbol.
func RandomCurrency() string {
	currencies := []string{"USD", "ARS", "EUR"}

	return currencies[rand.Intn(len(currencies))]
}

// RandomEmail random generates and email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
