package tool

import (
	"bytes"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// argon2id params
const (
	time       = 1
	memory     = 1024 * 64
	threads    = 2
	saltLength = 16
	keyLength  = 32
)

func Generate(password []byte) ([]byte, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate new salt: %w", err)
	}
	return append(salt, argon2.IDKey(password, salt, time, memory, threads, keyLength)...), nil
}

func Verify(password, hash []byte) (bool, error) {
	if len(hash) != saltLength+keyLength {
		return false, fmt.Errorf("invalid hash: length must be %d, got %d", saltLength+keyLength, len(hash))
	}
	if bytes.Equal(argon2.IDKey(password, hash[:saltLength], time, memory, threads, keyLength), hash[saltLength:]) {
		return true, nil
	}
	return false, nil
}
