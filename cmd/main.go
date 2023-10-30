package main

import (
	"log"
	"refactoring/internal/app/apiserver"
)

func main() {
	config := apiserver.NewConfig()

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
