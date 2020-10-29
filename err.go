package main

import (
	"errors"
	"log"
)

var (
	ErrRepoNotExist = errors.New("repository does not exist")
)

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}