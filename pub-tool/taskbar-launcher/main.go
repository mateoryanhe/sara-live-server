package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// Set at build time: go build -ldflags "-X main.scriptName=upload.bat"
var scriptName = "upload.bat"

func main() {
	exeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		waitExit(1)
		return
	}

	script := filepath.Join(exeDir, scriptName)
	if _, err := os.Stat(script); err != nil {
		fmt.Printf("Deploy script not found: %s\n", script)
		waitExit(1)
		return
	}

	cmd := exec.Command("cmd.exe", "/k", script)
	cmd.Dir = exeDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

func waitExit(code int) {
	fmt.Println("\nPress Enter to close...")
	fmt.Scanln()
	os.Exit(code)
}
