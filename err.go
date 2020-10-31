package main

import (
	"errors"
	"fmt"
	"log"
)

type Err struct {
	err error
	msg string
}

var (
	ErrNotPermitted = errors.New("Operation not permitted")
)

func (e *Err) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func errorWithMessage(err error, message string) *Err {
	return &Err{err: err, msg: message}
}

// Must is an alias for CheckErr
func Must(err error) {
	CheckErr(err)
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
