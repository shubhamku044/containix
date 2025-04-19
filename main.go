package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	cmd := exec.Command(filepath.Join("bin", "containix"))
	fmt.Println(cmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
