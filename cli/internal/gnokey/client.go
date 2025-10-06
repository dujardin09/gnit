package gnokey

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"gnit/cli/internal/config"
)

type Client struct {
	config *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{config: cfg}
}

func (c *Client) Query(expression string) ([]byte, error) {
	cmd := exec.Command("gnokey", "query", "vm/qeval",
		"-data", expression,
		"-remote", c.config.Remote)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return parseHexOutput(string(output))
}

func (c *Client) Run(gnoCode string) error {
	tmpFile := "/tmp/gnit_tx.gno"
	if err := os.WriteFile(tmpFile, []byte(gnoCode), 0644); err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile)

	cmd := exec.Command("gnokey", "maketx", "run",
		"-gas-fee", c.config.GasFee,
		"-gas-wanted", c.config.GasWanted,
		"-broadcast",
		"-chainid", c.config.ChainID,
		"-remote", c.config.Remote,
		c.config.Account,
		tmpFile)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
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
		return nil, fmt.Errorf("data: line not found in output")
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
		return nil, fmt.Errorf("hex decoding error: %w", err)
	}

	return data, nil
}
