package container

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

func randomHash() string {
	rand.Seed(time.Now().Unix())
	randBuffer := make([]byte, 32)
	rand.Read(randBuffer)
	sha := sha256.New().Sum(randBuffer)
	return fmt.Sprintf("%x", sha)[:64]
}
