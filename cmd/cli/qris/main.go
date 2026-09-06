package main

import (
	"os/exec"
)

func main() {
	exec.Command("powershell", "-Command", "Start-Process", "qris.jpg").Run()
}
