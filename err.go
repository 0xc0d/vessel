package main

import (
	"errors"
	"log"
)

var (
	ErrRepoNotExist = errors.New("repository does not exist")
)

// Must is an alias for CheckErr
func Must(err error) {
	CheckErr(err)
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}