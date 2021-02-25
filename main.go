package main

import (
	"log"
	"os"

	"github.com/Mayowa-Ojo/escobar/cmd"
)

func main() {
	app := cmd.Execute()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
