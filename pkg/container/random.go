package container

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func randomHash() string {
	randBuffer := make([]byte, 32)
	rand.Read(randBuffer)
	sha := sha256.New().Sum(randBuffer)
	return fmt.Sprintf("%x", sha)[:64]
}
