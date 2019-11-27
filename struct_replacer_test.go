package replacer_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/replacer"
)

func TestLoadStructReplacerFromFile(t *testing.T) {

	t.Run("should load pattern from file content", func(t *testing.T) {
		filecontents := `{
		"matchWith": "key",
		"replaceWith": "values",
		"patterns": {
			"pattern-1": "value-1",
			"pattern-2": "value-2"
		}
	}`

		expected := map[string]string{"pattern-1": "value-1", "pattern-2": "value-2"}

		reader := strings.NewReader(filecontents)

		r, err := replacer.LoadStructReplacerFromFile(reader)
		re, ok := r.(replacer.StructReplacer)

		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.True(t, ok)
		assert.NotNil(t, re)
		assert.Equal(t, "key", re.MatchWith)
		assert.Equal(t, "values", re.ReplaceWith)
		assert.Equal(t, expected, re.Patterns)
	})

	t.Run("should return error for invalid file content", func(t *testing.T) {
		invalidFilecontent := `{
		notajson
}`
		invalidJSONReader := strings.NewReader(invalidFilecontent)

		_, err := replacer.LoadStructReplacerFromFile(invalidJSONReader)

		assert.Error(t, err, "Should have failed parsing JSON")
	})
}

func TestNewStructReplacer(t *testing.T) {
	t.Run("should return new MapReplacer", func(t *testing.T) {
		patterns := map[string]string{"pattern": "value"}
		r := replacer.NewStructReplacer("name", "value", patterns)
		re, ok := r.(replacer.StructReplacer)

		assert.NotNil(t, r)
		assert.True(t, ok)
		assert.NotNil(t, re)
		assert.Equal(t, "name", re.MatchWith)
		assert.Equal(t, "value", re.ReplaceWith)
		assert.Equal(t, patterns, re.Patterns)
	})
}

func TestStructReplace(t *testing.T) {
	type entry struct {
		Name  string
		Value string
	}

	t.Run("When KeyField is different from ValueField", func(t *testing.T) {
		t.Run("should completely replace ValueField if pattern matches in KeyField", func(t *testing.T) {
			patterns := map[string]string{"pattern": "replaced value"}
			r := replacer.NewStructReplacer("Name", "Value", patterns)
			source := entry{Name: "pattern", Value: "value-1"}

			re, changed := r.Replace(&source)
			result, ok := re.(*entry)

			assert.NotNil(t, re)
			assert.True(t, changed)
			assert.True(t, ok)
			assert.Equal(t, "replaced value", result.Value)
		})

		t.Run("should completely replace ValueField if pattern partially matches in KeyField", func(t *testing.T) {
			patterns := map[string]string{"pattern": "replaced value"}
			r := replacer.NewStructReplacer("Name", "Value", patterns)
			source := entry{Name: "This is pattern-1", Value: "value-1"}

			re, changed := r.Replace(&source)
			result, ok := re.(*entry)

			assert.NotNil(t, re)
			assert.True(t, changed)
			assert.True(t, ok)
			assert.Equal(t, "replaced value", result.Value)
		})

		t.Run("should not replace ValueField if pattern does not match KeyField", func(t *testing.T) {
			patterns := map[string]string{"pattern": "replaced value"}
			r := replacer.NewStructReplacer("Name", "Value", patterns)
			source := entry{Name: "This is some key", Value: "value-1"}

			re, changed := r.Replace(&source)
			result, ok := re.(*entry)

			assert.NotNil(t, re)
			assert.False(t, changed)
			assert.True(t, ok)
			assert.Equal(t, "value-1", result.Value)
		})
	})

	t.Run("When KeyField and ValueField are same", func(t *testing.T) {
		t.Run("should replace value based on pattern matches", func(t *testing.T) {
			type config struct {
				Key   string
				Value string
			}
			type testCase struct {
				name          string
				input         config
				patterns      map[string]string
				expectedValue string
				isChanged     bool
			}

			testCases := []testCase{
				{
					name: "Matches with .* pattern",
					input: config{
						Key:   "Key-1",
						Value: "This is value.*",
					},
					patterns:      map[string]string{"value.*": "replaced value"},
					expectedValue: "This is replaced value",
					isChanged:     true,
				},
				{
					name: "Matches with .* pattern",
					input: config{
						Key:   "Key-1",
						Value: "This is value-1",
					},
					patterns:      map[string]string{"value.*": "replaced value"},
					expectedValue: "This is replaced value",
					isChanged:     true,
				},
				{
					name: ".default pattern",
					input: config{
						Key:   "Key-1",
						Value: "test.default",
					},
					patterns:      map[string]string{"test(.default)?$": "test-default.service"},
					expectedValue: "test-default.service",
					isChanged:     true,
				},
				{
					name: "simple service name match pattern",
					input: config{
						Key:   "Key-1",
						Value: "test",
					},
					patterns:      map[string]string{"test(.default)?$": "test-default.service"},
					expectedValue: "test-default.service",
					isChanged:     true,
				},
				{
					name: "No pattern match",
					input: config{
						Key:   "Key-1",
						Value: "test-db-core",
					},
					patterns:      map[string]string{"value.*": "replaced value", "test(.default)?$": "test-default.service"},
					expectedValue: "test-db-core",
					isChanged:     false,
				},
			}

			for _, test := range testCases {
				r := replacer.NewStructReplacer("Value", "Value", test.patterns)

				re, changed := r.Replace(&test.input)

				assert.NotNil(t, re, test.name)
				assert.Equal(t, test.isChanged, changed, test.name)

				result, ok := re.(*config)
				assert.True(t, ok, test.name)
				assert.Equal(t, test.expectedValue, result.Value, test.name)
			}
		})
	})

	t.Run("should not replace if the given input is not a pointer", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewStructReplacer("Value", "Value", patterns)
		source := entry{Name: "name", Value: "This is pattern-1"}

		re, changed := r.Replace(source)
		result, ok := re.(entry)

		assert.NotNil(t, re)
		assert.False(t, changed)
		assert.True(t, ok)
		assert.Equal(t, "This is pattern-1", result.Value)
	})

	t.Run("should not replace if the keyField doesn't match/exist", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewStructReplacer("DifferentKeyNotName", "Value", patterns)
		source := entry{Name: "pattern", Value: "value-1"}

		re, changed := r.Replace(&source)
		result, ok := re.(*entry)

		assert.NotNil(t, re)
		assert.False(t, changed)
		assert.True(t, ok)
		assert.Equal(t, "value-1", result.Value)
	})

	t.Run("should not replace if the valueField doesn't match/exist", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewStructReplacer("Name", "DifferentNotValue", patterns)
		source := entry{Name: "pattern", Value: "value-1"}

		re, changed := r.Replace(&source)
		result, ok := re.(*entry)

		assert.NotNil(t, re)
		assert.False(t, changed)
		assert.True(t, ok)
		assert.Equal(t, "value-1", result.Value)
	})

	t.Run("should not replace []items", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewStructReplacer("Name", "Value", patterns)
		source := []entry{{Name: "This is pattern-1", Value: "value-1"}}

		re, changed := r.Replace(&source)
		result, ok := re.(*[]entry)

		assert.NotNil(t, re)
		assert.False(t, changed)
		assert.True(t, ok)
		assert.Equal(t, entry{Name: "This is pattern-1", Value: "value-1"}, (*result)[0])
	})
}
