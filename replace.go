package replacer

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Replacer is the struct that holds the patterns for find and replace
type Replacer struct {
	patterns map[string]string
}

// New function creates a new Replacer
func New(p map[string]string) Replacer {
	return Replacer{patterns: p}
}

// LoadFromFile loads a map from a Reader
// The contents should be of the format
// {
//     "pattern-1": "replace-value-1",
//     "pattern-2": "replace-value-2"
// }
func LoadFromFile(r io.Reader) (Replacer, error) {
	var p map[string]string

	err := json.NewDecoder(r).Decode(&p)
	if err != nil {
		return Replacer{}, fmt.Errorf("Error loading from file: %v", err)
	}

	return New(p), nil
}

// Replace uses the patters configured in the replacer and accepts a string
// If a match is found in patterns, it replaces it with the value
// It also provides a bool, indicating if a change was done
func (r *Replacer) Replace(source string) (string, bool) {
	var changed bool
	for key, value := range r.patterns {
		if strings.Contains(source, key) {
			source = strings.Replace(source, key, value, -1)
			changed = true
		}
	}
	return source, changed
}
