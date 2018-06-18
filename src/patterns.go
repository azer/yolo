package yolo

import (
	"path/filepath"
)

type Patterns []string

func (patterns *Patterns) String() string {
	return "[patterns]"
}

func (patterns *Patterns) Set(value string) error {
	*patterns = append(*patterns, value)
	return nil
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
