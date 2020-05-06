package main

import (
	"log"

	"richard/rspace-client/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	operatorCommand := cmd.NewOperatorCommand()

	err := doc.GenMarkdownTree(operatorCommand, "./generated")
	if err != nil {
		log.Fatal(err)
	}
}
