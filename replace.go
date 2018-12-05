package replacer

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// MatcherField represents field type to be matched.
type MatcherField string

const (
	// VALUE represents matcher field as 'value'
	VALUE MatcherField = "value"
)

// Replacer is the struct that holds the patterns for find and replace
type Replacer struct {
	MatchWith MatcherField
	Patterns  Patterns
}

// Patterns represents field and value to be replaced with
type Patterns map[string]string

// New function creates a new Replacer
func New(p Patterns, with *string) Replacer {
	if with == nil {
		return Replacer{MatchWith: "value", Patterns: p}
	}
	return Replacer{MatchWith: MatcherField(*with), Patterns: p}
}

// LoadFromFile loads a map from a Reader
// The contents should be of the format
// {
//	"MatchWith": "value",
//	"Patterns": {
//     "pattern-1": "replace-value-1",
//     "pattern-2": "replace-value-2"
//	}
// }
func LoadFromFile(r io.Reader) (Replacer, error) {
	var rr Replacer

	if err := json.NewDecoder(r).Decode(&rr); err != nil {
		return Replacer{}, fmt.Errorf("error loading from file: %v", err)
	}

	if rr.MatchWith != VALUE {
		rr.MatchWith = VALUE
	}

	return rr, nil
}

// Replace uses the patters configured in the replacer and accepts a string
// If a match is found in patterns, it replaces it with the value
// It also provides a bool, indicating if a change was done
func (r *Replacer) Replace(source string) (string, bool) {
	var changed bool
	for key, value := range r.Patterns {
		if strings.Contains(source, key) {
			source = strings.Replace(source, key, value, -1)
			changed = true
		}
	}
	return source, changed
}
