package replacer

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// StringReplacer is the struct that holds the patterns for find and replace
type StringReplacer struct {
	Patterns map[string]string
}

// NewStringReplacer function creates a new Replacer
func NewStringReplacer(p map[string]string) Replacer {
	return StringReplacer{Patterns: p}
}

// LoadStringReplacerFromFile loads a map from a Reader
// The contents should be of the format
// {
//		"patterns": {
//			"pattern-1": "replace-value-1",
//			"pattern-2": "replace-value-2"
//		}
// }
func LoadStringReplacerFromFile(r io.Reader) (Replacer, error) {
	replacer := StringReplacer{}

	err := json.NewDecoder(r).Decode(&replacer)
	if err != nil {
		return nil, fmt.Errorf("Error loading from file: %v", err)
	}

	return replacer, nil
}

// Replace uses the patters configured in the replacer and accepts a string
// If a match is found in patterns, it replaces it with the value
// It also provides a bool, indicating if a change was done
func (r StringReplacer) Replace(source interface{}) (interface{}, bool) {
	if source, ok := source.(string); ok == true {
		var changed bool
		for key, value := range r.Patterns {
			if strings.Contains(source, key) {
				source = strings.Replace(source, key, value, -1)
				changed = true
			}
		}
		return source, changed
	}
	return source, false
}
