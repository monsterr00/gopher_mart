package main

import (
	"log"

	"github.com/monsterr00/gopher_mart/internal"
)

func main() {
	err := internal.Run()

	if err != nil {
		log.Fatal(err)
	}
}
