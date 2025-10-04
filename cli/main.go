package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "pull":
		handlePull()
	case "commit":
		handleCommit()
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Error: unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: gnit <command> [options]")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  pull <file>     Fetch a file from the repository")
	fmt.Println("  commit <message>  Commit changes with a message")
	fmt.Println("  help            Display this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gnit pull example.gno")
	fmt.Println("  gnit commit \"My commit message\"")
}

func handlePull() {
	if len(os.Args) < 3 {
		fmt.Println("Error: filename required for pull")
		fmt.Println("Usage: gnit pull <file>")
		os.Exit(1)
	}

	filename := os.Args[2]

	fmt.Printf("Pulling '%s'...\n", filename)

	cmd := exec.Command("gnokey", "query", "vm/qeval",
		"-data", fmt.Sprintf("gno.land/r/example.Pull(\"%s\")", filename),
		"-remote", "tcp://127.0.0.1:26657")

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}

	content, err := parseHexOutput(string(output))
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		os.Exit(1)
	}

	if len(content) == 0 {
		fmt.Printf("File '%s' not found or empty\n", filename)
		os.Exit(1)
	}

	err = os.WriteFile(filename, content, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("File '%s' fetched successfully (%d bytes)\n", filename, len(content))
}

func handleCommit() {
	if len(os.Args) < 3 {
		fmt.Println("Error: message required for commit")
		fmt.Println("Usage: gnit commit \"<message>\"")
		os.Exit(1)
	}

	message := strings.Join(os.Args[2:], " ")
	message = strings.Trim(message, "\"")

	fmt.Printf("Committing with message: '%s'...\n", message)

	files := make(map[string][]byte)

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				fmt.Printf("Warning: unable to read %s: %v\n", path, readErr)
				return nil
			}
			files[path] = content
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error reading files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No .md files found to commit")
		os.Exit(1)
	}

	fmt.Printf("Files to commit: %d\n", len(files))
	for filename := range files {
		fmt.Printf("  - %s\n", filename)
	}

	var filesData strings.Builder
	for filename, content := range files {
		filesData.WriteString(filename)
		filesData.WriteString("|")
		filesData.WriteString(string(content))
		filesData.WriteString("\n")
	}

	cmd := exec.Command("gnokey", "maketx", "call",
		"-pkgpath", "gno.land/r/example",
		"-func", "Push",
		"-args", message,
		"-args", filesData.String(),
		"-gas-fee", "1000000ugnot",
		"-gas-wanted", "50000000",
		"-send", "",
		"-broadcast",
		"-chainid", "dev",
		"-remote", "tcp://127.0.0.1:26657",
		"test")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Commit successful!\n")
}

func parseHexOutput(output string) ([]byte, error) {
	lines := strings.Split(output, "\n")
	var dataLine string
	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			dataLine = line
			break
		}
	}

	if dataLine == "" {
		return nil, fmt.Errorf("data: line not found")
	}

	if strings.Contains(dataLine, "(nil []uint8)") {
		return []byte{}, nil
	}

	re := regexp.MustCompile(`slice\[0x([0-9a-fA-F]+)\]`)
	matches := re.FindStringSubmatch(dataLine)

	if len(matches) < 2 {
		return nil, fmt.Errorf("unrecognized output format: %s", dataLine)
	}

	data, err := hex.DecodeString(matches[1])
	if err != nil {
		return nil, fmt.Errorf("hex decoding error: %v", err)
	}

	return data, nil
}
