package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}
	if nRead < n {
		return nil, fmt.Errorf("bytes: didnt read enough random bytes")
	}
	return b, nil
}

// n is the number of bytes being used to generate a random string.
func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("String: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

const SessionTokenBytes = 32

func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}
