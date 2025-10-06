package main

import (
	"fmt"
	"os"
	"strings"

	"gnit/cli/internal/commands"
	"gnit/cli/internal/config"
	"gnit/cli/internal/gnokey"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg := config.DefaultConfig()

	client := gnokey.NewClient(cfg)

	command := os.Args[1]

	switch command {
	case "pull":
		handlePull(client, cfg)
	case "commit":
		handleCommit(client, cfg)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Error: unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handlePull(client *gnokey.Client, cfg *config.Config) {
	if len(os.Args) < 3 {
		fmt.Println("Error: filename required for pull")
		fmt.Println("Usage: gnit pull <file>")
		os.Exit(1)
	}

	filename := os.Args[2]
	cmd := commands.NewPull(client, cfg)

	if err := cmd.Execute(filename); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func handleCommit(client *gnokey.Client, cfg *config.Config) {
	if len(os.Args) < 3 {
		fmt.Println("Error: message required for commit")
		fmt.Println("Usage: gnit commit \"<message>\"")
		os.Exit(1)
	}

	message := strings.Join(os.Args[2:], " ")
	message = strings.Trim(message, "\"")

	cmd := commands.NewCommit(client, cfg)

	if err := cmd.Execute(message); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: gnit <command> [options]")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  pull <file>       Fetch a file from the repository")
	fmt.Println("  commit <message>  Commit changes with a message")
	fmt.Println("  help              Display this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gnit pull example.gno")
	fmt.Println("  gnit commit \"My commit message\"")
}
