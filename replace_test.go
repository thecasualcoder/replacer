package replacer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromFile(t *testing.T) {
	t.Run("Should match patterns with value", func(t *testing.T) {
		filecontents := `{
		"MatchWith": "value",
		"Patterns": {
			"pattern-1": "value-1",
			"pattern-2": "value-2"
		}
	}`
		reader := strings.NewReader(filecontents)

		r, err := LoadFromFile(reader)
		expected := Replacer{
			MatchWith: "value",
			Patterns: Patterns{
				"pattern-1": "value-1",
				"pattern-2": "value-2",
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, r)

		invalidFilecontent := `{
		notajson
}`
		invalidJSONReader := strings.NewReader(invalidFilecontent)

		_, err = LoadFromFile(invalidJSONReader)

		assert.Error(t, err, "Should have failed parsing JSON")
	})

	t.Run("Should match patterns with value even if MatchWith is invalid", func(t *testing.T) {
		filecontents := `{
		"MatchWith": "breakMe",
		"Patterns": {
			"pattern-1": "value-1",
			"pattern-2": "value-2"
		}
	}`
		reader := strings.NewReader(filecontents)

		r, err := LoadFromFile(reader)
		expected := Replacer{
			MatchWith: "value",
			Patterns: Patterns{
				"pattern-1": "value-1",
				"pattern-2": "value-2",
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, r)

		invalidFilecontent := `{
		notajson
}`
		invalidJSONReader := strings.NewReader(invalidFilecontent)

		_, err = LoadFromFile(invalidJSONReader)

		assert.Error(t, err, "Should have failed parsing JSON")
	})
}

func TestReplace(t *testing.T) {
	t.Run("Should replace values based on patterns", func(t *testing.T) {
		patterns := Patterns{
			"pattern-1": "value-1",
			"pattern-2": "value-2",
		}

		r := New(patterns, nil)

		result, changed := r.Replace("This is pattern-1")

		assert.True(t, changed)
		assert.Equal(t, "This is value-1", result)

		result, changed = r.Replace("This need not be changed")

		assert.False(t, changed)
		assert.Equal(t, "This need not be changed", result)

		result, changed = r.Replace("This is pattern-1 and this is pattern-2")

		assert.True(t, changed)
		assert.Equal(t, "This is value-1 and this is value-2", result)
	})
}
