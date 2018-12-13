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
		"keyField": "key",
		"valueField": "values",
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
		assert.Equal(t, "key", re.KeyField)
		assert.Equal(t, "values", re.ValueField)
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
		assert.Equal(t, "name", re.KeyField)
		assert.Equal(t, "value", re.ValueField)
		assert.Equal(t, patterns, re.Patterns)
	})
}

func TestStructReplace(t *testing.T) {
	type example struct {
		Name  string
		Value string
	}

	t.Run("should completely replace value FieldKey is different from FieldValue and pattern matches", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewStructReplacer("Name", "Value", patterns)
		source := example{Name: "This is pattern-1", Value: "value-1"}

		re, changed := r.Replace(&source)
		result, ok := re.(*example)

		assert.NotNil(t, re)
		assert.True(t, changed)
		assert.True(t, ok)
		assert.Equal(t, "replaced value", result.Value)
	})

	t.Run("should do pattern based replace if KeyField is different from ValueField and pattern matches", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewStructReplacer("Value", "Value", patterns)
		source := example{Name: "name", Value: "This is pattern-1"}

		re, changed := r.Replace(&source)
		result, ok := re.(*example)

		assert.NotNil(t, re)
		assert.True(t, changed)
		assert.True(t, ok)
		assert.Equal(t, "This is replaced value-1", result.Value)
	})

	t.Run("should not replace if the value source is not ptr", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewStructReplacer("Value", "Value", patterns)
		source := example{Name: "name", Value: "This is pattern-1"}

		re, changed := r.Replace(source)
		result, ok := re.(example)

		assert.NotNil(t, re)
		assert.False(t, changed)
		assert.True(t, ok)
		assert.Equal(t, "This is pattern-1", result.Value)
	})

	t.Run("should not replace if the keyField doesn't match", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewStructReplacer("Name", "Value", patterns)
		source := example{Name: "key", Value: "value-1"}

		re, changed := r.Replace(&source)
		result, ok := re.(*example)

		assert.NotNil(t, re)
		assert.False(t, changed)
		assert.True(t, ok)
		assert.Equal(t, "value-1", result.Value)
	})

	t.Run("should not replace if the keyField doesn't exist", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewStructReplacer("Key", "Value", patterns)
		source := example{Name: "This is pattern-1", Value: "value-1"}

		re, changed := r.Replace(&source)
		result, ok := re.(*example)

		assert.NotNil(t, re)
		assert.False(t, changed)
		assert.True(t, ok)
		assert.Equal(t, "value-1", result.Value)
	})

	t.Run("should not replace if the valueField doesn't match", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewStructReplacer("Value", "Value", patterns)
		source := example{Name: "key", Value: "value-1"}

		re, changed := r.Replace(&source)
		result, ok := re.(*example)

		assert.NotNil(t, re)
		assert.False(t, changed)
		assert.True(t, ok)
		assert.Equal(t, "value-1", result.Value)
	})

	t.Run("should not replace if the valueField doesn't exist", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewStructReplacer("Name", "Age", patterns)
		source := example{Name: "This is pattern-1", Value: "value-1"}

		re, changed := r.Replace(&source)
		result, ok := re.(*example)

		assert.NotNil(t, re)
		assert.False(t, changed)
		assert.True(t, ok)
		assert.Equal(t, "value-1", result.Value)
	})

	t.Run("should not replace []items", func(t *testing.T) {
		patterns := map[string]string{"pattern": "replaced value"}
		r := replacer.NewStructReplacer("Name", "Value", patterns)
		source := []example{{Name: "This is pattern-1", Value: "value-1"}}

		re, changed := r.Replace(&source)
		result, ok := re.(*[]example)

		assert.NotNil(t, re)
		assert.False(t, changed)
		assert.True(t, ok)
		assert.Equal(t, example{Name: "This is pattern-1", Value: "value-1"}, (*result)[0])
	})
}
