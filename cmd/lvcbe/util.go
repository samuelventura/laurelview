package main

import (
	"crypto/sha256"
	"fmt"
)

func hashit(data string) string {
	sum := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", sum)
}

func pwdit(id string) string {
	return hashit(id)[:8]
}
