package commands

import (
	"fmt"

	"gnit/cli/internal/config"
	"gnit/cli/internal/filesystem"
	"gnit/cli/internal/gnokey"
)

type Pull struct {
	client *gnokey.Client
	config *config.Config
}

func NewPull(client *gnokey.Client, cfg *config.Config) *Pull {
	return &Pull{
		client: client,
		config: cfg,
	}
}

func (p *Pull) Execute(filename string) error {
	fmt.Printf("Pulling '%s'...\n", filename)

	query := fmt.Sprintf("%s.Repo.Pull(\"%s\")", p.config.RealmPath, filename)
	content, err := p.client.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query file: %w", err)
	}

	if len(content) == 0 {
		return fmt.Errorf("file '%s' not found or empty", filename)
	}

	if err := filesystem.WriteFile(filename, content); err != nil {
		return err
	}

	fmt.Printf("File '%s' fetched successfully (%d bytes)\n", filename, len(content))
	return nil
}
