package main

import (
	"fmt"
	"path/filepath"

	"github.com/umardev500/jnt-report/internal/auto"
)

func main() {

	path := filepath.Join(
		"C:\\", "Users", "User", "Projects", "go-report",
		"generated", "excel", "sj", "2026_05_04",
		"sj_print_20260504_231111.xlsm",
	)

	fmt.Println(path)
	auto.PrintSJ(path)
	fmt.Println("done")
}
