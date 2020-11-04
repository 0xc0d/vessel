package container

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"strings"
)

// randomHash creates random Length 64 hash
func randomHash() string {
	randBuffer := make([]byte, 32)
	rand.Read(randBuffer)
	sha := sha256.New().Sum(randBuffer)
	return fmt.Sprintf("%x", sha)[:64]
}

// completeDigest completes a prefix to a length 64 digest.
func completeDigest(prefix string) (digest string) {
	if len(prefix) == DigestStdLen {
		return prefix
	}
	list, err := ioutil.ReadDir(containerPath)
	if err != nil {
		return
	}
	for _, file := range list {
		if file.IsDir() && strings.HasPrefix(file.Name(), prefix) {
			return file.Name()
		}
	}
	return
}