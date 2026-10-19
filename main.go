package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	var inputData []byte
	var err error

	// Inspect stdin mode to check for incoming pipes or redirects
	stat, statErr := os.Stdin.Stat()
	isPiped := (statErr == nil) && ((stat.Mode() & os.ModeCharDevice) == 0)

	args := os.Args[1:]

	if len(args) > 0 {
		// Read and concatenate files if arguments are supplied
		var buf bytes.Buffer
		for _, path := range args {
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				fmt.Fprintf(os.Stderr, "copey: failed to read %s: %v\n", path, readErr)
				os.Exit(1)
			}
			buf.Write(data)
		}
		inputData = buf.Bytes()
	} else if isPiped {
		// Read directly from stdin stream
		inputData, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "copey: failed to read from stdin: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Usage: copey <file> ... or <command> | copey\n")
		os.Exit(1)
	}

	if len(inputData) == 0 {
		os.Exit(0)
	}

	if err := copyToClipboard(inputData); err != nil {
		fmt.Fprintf(os.Stderr, "copey: error copying to clipboard: %v\n", err)
		os.Exit(1)
	}
}

// copyToClipboard writes byte data directly into the platform's native clipboard tool.
func copyToClipboard(data []byte) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		if os.Getenv("WAYLAND_DISPLAY") != "" && commandExists("wl-copy") {
			cmd = exec.Command("wl-copy")
		} else if commandExists("xclip") {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if commandExists("xsel") {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("no clipboard utility found (install wl-clipboard, xclip, or xsel)")
		}
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	cmd.Stdin = bytes.NewReader(data)
	return cmd.Run()
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
