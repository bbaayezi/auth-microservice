package utils

import (
	"math/rand"
	"time"
)

// This file provides functions to generate random string of fixed length
const (
	charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func nonceWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateNonce generates random string with fixed length
func GenerateNonce(length int) string {
	return nonceWithCharset(length, charset)
}
