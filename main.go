package main

import (
	"crud/cmd"
	"log"
)

func main() {
	err := cmd.RunServer()
	if err != nil {
		log.Fatal(err)
	}
}
