package replacer

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"regexp"
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
	var changed bool
	srcValue := reflect.ValueOf(source)
	v := srcValue.Elem()
	matchField := v.FieldByName(r.MatchWith)
	replaceField := v.FieldByName(r.ReplaceWith)

	if matchField.IsValid() &&
		matchField.Kind() == reflect.String &&
		replaceField.IsValid() &&
		replaceField.CanSet() &&
		replaceField.Kind() == reflect.String {
		for regexPattern, toBeReplacedWith := range r.Patterns {
			var re = regexp.MustCompile(regexPattern)
			finalValue := toBeReplacedWith
			if re.MatchString(matchField.String()) {
				if r.MatchWith == r.ReplaceWith {
					finalValue = re.ReplaceAllString(replaceField.String(), toBeReplacedWith)
				}
				replaceField.SetString(finalValue)
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
