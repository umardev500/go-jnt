package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exePath)

	fmt.Println("Binary path:", exePath)
	fmt.Println("Binary directory:", exeDir)
}
