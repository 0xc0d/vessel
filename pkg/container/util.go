package container

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"strings"
)

func randomHash() string {
	randBuffer := make([]byte, 32)
	rand.Read(randBuffer)
	sha := sha256.New().Sum(randBuffer)
	return fmt.Sprintf("%x", sha)[:64]
}

func completeDigest(prefix string) (digest string) {
	if len(prefix) == DigestStdLen {
		return prefix
	}
	list, err := ioutil.ReadDir(CtrDir)
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