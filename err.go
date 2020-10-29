package main

import (
	"log"
)

func must(err error)  {
	if err != nil {
		log.Fatal(err)
	}
}