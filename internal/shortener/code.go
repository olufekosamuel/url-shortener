package shortener

import (
	"crypto/rand"
	"errors"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandomCode returns n random characters from the Base62 alphabet.
//
// Base62 stays URL-safe (no +, /, or punctuation). Seven characters give
// 62^7 possible codes, so collisions are rare; Save still rejects duplicates
// and Shorten retries.
func RandomCode(n int) (string, error) {
	if n <= 0 {
		return "", errors.New("code length must be positive")
	}

	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	out := make([]byte, n)
	for i, b := range buf {
		out[i] = alphabet[int(b)%len(alphabet)]
	}
	return string(out), nil
}
