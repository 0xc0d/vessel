package main

import (
	"errors"
	"log"
)

var (
	ErrRepoNotExist = errors.New("repository does not exist")
)

func must(err error)  {
	if err != nil {
		log.Fatal(err)
	}
}