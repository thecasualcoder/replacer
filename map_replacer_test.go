package replacer_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/replacer"
)

func TestLoadMapReplacerFromFile(t *testing.T) {
	t.Run("should load pattern from file content", func(t *testing.T) {
		filecontents := `{
		"matchWith": "values",
		"patterns": {
			"pattern-1": "value-1",
			"pattern-2": "value-2"
		}
	}`

		expected := map[string]string{"pattern-1": "value-1", "pattern-2": "value-2"}

		reader := strings.NewReader(filecontents)

		r, err := replacer.LoadMapReplacerFromFile(reader)
		re, ok := r.(replacer.MapReplacer)

		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.True(t, ok)
		assert.NotNil(t, re)
		assert.Equal(t, "values", re.MatchWith)
		assert.Equal(t, expected, re.Patterns)
	})

	t.Run("should return error for invalid file content", func(t *testing.T) {
		invalidFilecontent := `{
		notajson
}`
		invalidJSONReader := strings.NewReader(invalidFilecontent)

		_, err := replacer.LoadMapReplacerFromFile(invalidJSONReader)

		assert.Error(t, err, "Should have failed parsing JSON")
	})
}

func TestNewMapReplacer(t *testing.T) {
	t.Run("should return new MapReplacer", func(t *testing.T) {
		patterns := map[string]string{"pattern": "value"}
		r := replacer.NewMapReplacer("values", patterns)
		re, ok := r.(replacer.MapReplacer)

		assert.NotNil(t, r)
		assert.True(t, ok)
		assert.NotNil(t, re)
		assert.Equal(t, "values", re.MatchWith)
		assert.Equal(t, patterns, re.Patterns)
	})
}

func TestMapReplace(t *testing.T) {
	t.Run("should replace value if key pattern matches", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewMapReplacer("key", patterns)

		result, changed := r.Replace(map[string]string{"This is pattern-1": "value-1", "something": "random"})
		expected := map[string]string{"This is pattern-1": "replaced value", "something": "random"}
		assert.True(t, changed)
		assert.Equal(t, expected, result)
	})

	t.Run("should replace value if value pattern matches", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewMapReplacer("value", patterns)

		result, changed := r.Replace(map[string]string{"key-1": "This is pattern-1", "key-2": "random"})
		expected := map[string]string{"key-1": "This is replaced value-1", "key-2": "random"}
		assert.True(t, changed)
		assert.Equal(t, expected, result)
	})

	t.Run("should not replace value if source is not map[string]string", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewMapReplacer("key", patterns)

		result, changed := r.Replace("source")
		assert.False(t, changed)
		assert.Equal(t, "source", result)
	})
}
