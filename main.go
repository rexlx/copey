package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

// copyToClipboard copies the given string to the system clipboard
func copyToClipboard(content string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("pbcopy")
	case "linux": // Linux
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows": // Windows
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Get the stdin pipe to write our content
	in, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	// Write the content to the pipe
	_, err = io.WriteString(in, content)
	if err != nil {
		return fmt.Errorf("error writing to pipe: %v", err)
	}

	// Close the pipe
	if err := in.Close(); err != nil {
		return fmt.Errorf("error closing pipe: %v", err)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command finished with error: %v", err)
	}

	return nil
}

func main() {
	// Define command-line flag for the file path
	filePath := flag.String("file", "", "Path to the file to copy to clipboard")
	flag.Parse()

	// Check if file path was provided
	if *filePath == "" {
		fmt.Println("Error: Please provide a file path using -file flag")
		fmt.Println("Usage: program -file <path_to_file>")
		os.Exit(1)
	}

	// Open the file
	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read the file contents
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Copy to clipboard
	err = copyToClipboard(string(content))
	if err != nil {
		fmt.Printf("Error copying to clipboard: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully copied contents of %s to clipboard\n", *filePath)
}
