package commands

import (
	"fmt"

	"gnit/cli/internal/config"
	"gnit/cli/internal/filesystem"
	"gnit/cli/internal/gnokey"
)

type Commit struct {
	client *gnokey.Client
	config *config.Config
}

func NewCommit(client *gnokey.Client, cfg *config.Config) *Commit {
	return &Commit{
		client: client,
		config: cfg,
	}
}

func (c *Commit) Execute(message string) error {
	fmt.Printf("Committing with message: '%s'...\n", message)

	files, err := filesystem.CollectFiles(".md")
	if err != nil {
		return fmt.Errorf("failed to collect files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no .md files found to commit")
	}

	fmt.Printf("Files to commit: %d\n", len(files))
	for filename := range files {
		fmt.Printf("  - %s\n", filename)
	}

	filesData := filesystem.SerializeFiles(files)

	gnoCode := c.generateCommitCode(message, filesData)

	if err := c.client.Run(gnoCode); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	fmt.Printf("Commit successful!\n")
	return nil
}

func (c *Commit) generateCommitCode(message, filesData string) string {
	return fmt.Sprintf(`package main

import (
	"strings"
	"gno.land/r/example"
)

func unescape(s string) string {
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\|", "|")
	s = strings.ReplaceAll(s, "\\\\", "\\")
	return s
}

func main() {
	filesData := %q
	lines := strings.Split(filesData, "\n")
	files := make(map[string][]byte)

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) == 2 {
			files[parts[0]] = []byte(unescape(parts[1]))
		}
	}

	hash := example.Repo.Commit(%q, files)
	println("Commit hash:", hash)
}
`, filesData, message)
}
