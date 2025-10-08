package ignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Matcher struct {
	patterns []string
}

func NewMatcher(root string) (*Matcher, error) {
	m := &Matcher{
		patterns: []string{},
	}

	ignorePath := filepath.Join(root, ".gnitignore")
	file, err := os.Open(ignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return m, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		m.patterns = append(m.patterns, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Matcher) Match(path string) bool {
	// Only ignore .git directory, not files starting with .git
	if strings.Contains(path, ".git/") || path == ".git" {
		return true
	}

	for _, pattern := range m.patterns {
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && matched {
			return true
		}

		if strings.HasPrefix(pattern, "*") {
			suffix := strings.TrimPrefix(pattern, "*")
			if strings.HasSuffix(path, suffix) {
				return true
			}
		} else if strings.Contains(path, pattern) {
			return true
		}
	}

	return false
}
