package replacer_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/replacer"
)

func TestLoadStringReplacerFromFile(t *testing.T) {

	t.Run("should load patterns", func(t *testing.T) {
		filecontents := `{
		"pattern-1": "value-1",
		"pattern-2": "value-2"
}`
		reader := strings.NewReader(filecontents)

		r, err := replacer.LoadStringReplacerFromFile(reader)
		re, ok := r.(replacer.StringReplacer)
		expected := map[string]string{
			"pattern-1": "value-1",
			"pattern-2": "value-2",
		}

		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NotNil(t, re)
		assert.Equal(t, expected, re.Patterns)
	})

	t.Run("should return error for invalid file content", func(t *testing.T) {
		invalidFilecontent := `{
		notajson
}`
		invalidJSONReader := strings.NewReader(invalidFilecontent)

		_, err := replacer.LoadStringReplacerFromFile(invalidJSONReader)

		assert.Error(t, err, "Should have failed parsing JSON")
	})
}

func TestStringReplace(t *testing.T) {
	t.Run("should replace source if its string", func(t *testing.T) {
		patterns := map[string]string{
			"pattern-1": "value-1",
			"pattern-2": "value-2",
		}

		r := replacer.NewStringReplacer(patterns)

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

	t.Run("should not replace source if its not string", func(t *testing.T) {
		patterns := map[string]string{}

		r := replacer.NewStringReplacer(patterns)

		result, changed := r.Replace(1)

		assert.False(t, changed)
		assert.Equal(t, 1, result)
	})
}
