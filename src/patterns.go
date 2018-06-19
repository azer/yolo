package yolo

import (
	"path/filepath"
	"strings"
)

type Patterns []string

func (patterns *Patterns) String() string {
	return "[patterns]"
}

func (patterns *Patterns) Set(value string) error {
	*patterns = append(*patterns, value)
	return nil
}

func (patterns Patterns) Has(filename string) bool {
	if strings.HasPrefix(filename, "./") {
		filename = filename[2:]
	}

	for _, pattern := range patterns {
		if pattern == filename || strings.HasPrefix(filename, pattern+"/") {
			return true
		}
	}

	return false
}

func (patterns Patterns) Expand() (*Patterns, error) {
	result := Patterns{}

	for _, pattern := range patterns {
		match, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}

		result = append(result, match...)
	}

	return &result, nil
}
