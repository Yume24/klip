package main

import (
	"fmt"
	"klip/internal/cli"
	"os"
)

const name = "Klip"

func main() {
	parser := cli.NewParser(name, os.Args[1:]) // First arg is the command itself
	config, err := parser.ParseArguments()

	if err != nil {
		fmt.Println(err)
		parser.FlagSet.Usage()
		os.Exit(2)
	}

	fmt.Println(config)
}