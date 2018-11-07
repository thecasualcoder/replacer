package replacer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromFile(t *testing.T) {
	filecontents := `{
		"pattern-1": "value-1",
		"pattern-2": "value-2"
}`
	reader := strings.NewReader(filecontents)

	r, err := LoadFromFile(reader)
	expected := map[string]string{
		"pattern-1": "value-1",
		"pattern-2": "value-2",
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, r.patterns)

	invalidFilecontent := `{
		notajson
}`
	invalidJsonReader := strings.NewReader(invalidFilecontent)

	_, err = LoadFromFile(invalidJsonReader)

	assert.Error(t, err, "Should have failed parsing JSON")
}

func TestReplace(t *testing.T) {
	patterns := map[string]string{
		"pattern-1": "value-1",
		"pattern-2": "value-2",
	}

	r := New(patterns)

	result, changed := r.Replace("This is pattern-1")

	assert.True(t, changed)
	assert.Equal(t, "This is value-1", result)

	result, changed = r.Replace("This need not be changed")

	assert.False(t, changed)
	assert.Equal(t, "This need not be changed", result)
}
