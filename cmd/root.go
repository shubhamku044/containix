package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "hello":
		fmt.Println("hello world")
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printHelp()
	}
}

func printHelp() {
	fmt.Println(`Containix - Docker Management CLI
Usage:
  containix <command>

Available Commands:
  hello        Print 'hello world'

Example:
  containix hello`)
}
