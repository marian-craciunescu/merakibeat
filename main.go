package main

import (
	"fmt"
	"github.com/marian-craciunescu/merakibeat/cmd"
	"log"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println("Error=", err.Error())
		log.Fatal("Could not start merakibeat")
	}
}
