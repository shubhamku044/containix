package main

import (
	"fmt"

	"github.com/shubhamku044/containix/internal/app"
)

func main() {
	// Run the application directly without building a binary
	fmt.Println("Starting Containix...")

	// Call the application's Run function directly
	app.Run()
}
