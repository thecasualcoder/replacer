package replacer

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// StructReplacer is the struct that holds the patterns for find and replace
type StructReplacer struct {
	MatchWith   string
	ReplaceWith string
	Patterns    map[string]string
}

// NewStructReplacer function creats and returns MapReplacer
func NewStructReplacer(matchWith string, replaceWith string, patterns map[string]string) Replacer {
	return StructReplacer{MatchWith: matchWith, ReplaceWith: replaceWith, Patterns: patterns}
}

// LoadStructReplacerFromFile load a map from a Reader
// The contents should be of the format
// {
//		"matchWith": "key",
//		"replaceWith": "value",
//		"patterns: {
//			"pattern-1": "replace-value-1",
//			"pattern-2": "replace-value-2"
//		}
// }
func LoadStructReplacerFromFile(r io.Reader) (Replacer, error) {
	replacer := StructReplacer{}

	if err := json.NewDecoder(r).Decode(&replacer); err != nil {
		return nil, fmt.Errorf("Error loading from file %v", err)
	}
	return replacer, nil
}

func (r StructReplacer) replaceStruct(source interface{}) (interface{}, bool) {
	srcValue := reflect.ValueOf(source)
	var changed bool
	v := srcValue.Elem()
	mf := v.FieldByName(r.MatchWith)
	rf := v.FieldByName(r.ReplaceWith)

	if mf.IsValid() && mf.Kind() == reflect.String && rf.IsValid() && rf.CanSet() && rf.Kind() == reflect.String {
		for key, value := range r.Patterns {
			if strings.Contains(mf.String(), key) {
				if r.MatchWith == r.ReplaceWith {
					value = strings.Replace(rf.String(), key, value, -1)
				}
				rf.SetString(value)
				changed = true
			}
		}
		return source, changed
	}
	return source, false
}

// Replace uses the patters configured in the replacer and accepts a struct
// If a match is found in patterns, it replaces it with the value
// It also provides a bool, indicating if a change was done
func (r StructReplacer) Replace(source interface{}) (interface{}, bool) {
	v := reflect.ValueOf(source)
	if v.Kind() != reflect.Ptr {
		return source, false
	}
	srcValue := reflect.TypeOf(source).Elem()
	switch srcValue.Kind() {
	case reflect.Struct:
		return r.replaceStruct(source)
	default:
		return source, false
	}
}
