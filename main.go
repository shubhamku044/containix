package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// This is a wrapper that calls the actual application in cmd/containix
	// This allows for backward compatibility while using the new structure
	cmd := exec.Command(filepath.Join("cmd", "containix", "containix"))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
