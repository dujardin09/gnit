package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CollectFiles(pattern string) (map[string][]byte, error) {
	files := make(map[string][]byte)

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, pattern) {
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
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return files, nil
}

func WriteFile(path string, content []byte) error {
	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}
	return nil
}

func SerializeFiles(files map[string][]byte) string {
	var builder strings.Builder
	for filename, content := range files {
		builder.WriteString(filename)
		builder.WriteString("|")
		builder.WriteString(string(content))
		builder.WriteString("\n")
	}
	return builder.String()
}
