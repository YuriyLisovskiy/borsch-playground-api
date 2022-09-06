package main

import (
	"log"

	"github.com/YuriyLisovskiy/borsch-playground-service/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Fatal(err)
	}
}
