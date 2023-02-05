package flatthis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
	testMap := map[string]any{
		"a": "a",
		"nest": map[string]any{
			"b": "b",
		},
	}
	f := New()
	flatted := f.FlatMap(testMap)
	assert.Equal(t, "a", flatted["a"])
	assert.Equal(t, "b", flatted["nest.b"])

}
