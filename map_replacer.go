package replacer

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// MapReplacer is the struct that holds the patterns for find and replace
type MapReplacer struct {
	MatchWith string
	Patterns  map[string]string
}

// NewMapReplacer function creats and returns MapReplacer
func NewMapReplacer(with string, patterns map[string]string) Replacer {
	return MapReplacer{MatchWith: with, Patterns: patterns}
}

// LoadMapReplacerFromFile load a map from a Reader
// The contents should be of the format
// {
//		"matchWith": "key",
//		"patterns: {
//			"pattern-1": "replace-value-1",
//			"pattern-2": "replace-value-2"
//		}
// }
func LoadMapReplacerFromFile(r io.Reader) (Replacer, error) {
	replacer := MapReplacer{}

	if err := json.NewDecoder(r).Decode(&replacer); err != nil {
		return nil, fmt.Errorf("Error loading from file %v", err)
	}
	return replacer, nil
}

// Replace uses the patters configured in the replacer and accepts a map[string]interface{}
// If a match is found in patterns, it replaces it with the value
// It also provides a bool, indicating if a change was done
func (r MapReplacer) Replace(source interface{}) (interface{}, bool) {
	if source, ok := source.(map[string]string); ok {
		result := map[string]string{}
		var changed bool
		for pattern, valueToBeReplaced := range r.Patterns {
			for key, value := range source {
				switch r.MatchWith {
				case "key":
					if strings.Contains(key, pattern) {
						result[key] = valueToBeReplaced
						changed = true
						continue
					}
				default:
					if strings.Contains(value, pattern) {
						result[key] = strings.Replace(value, pattern, valueToBeReplaced, -1)
						changed = true
						continue
					}
				}
				result[key] = value
			}
		}
		return result, changed
	}
	return source, false
}
