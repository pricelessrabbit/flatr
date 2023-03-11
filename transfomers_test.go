package flatr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUseFieldAsIndex(t *testing.T) {
	toTest := map[string]any{
		"foo": []any{
			map[string]any{"id": "a", "data": 10},
			map[string]any{"id": "b", "data": 20},
		},
		"bar": []any{
			map[string]any{"id": 1, "data": 30},
			map[string]any{"id": 2, "data": 40},
		},
	}
	f := New(
		AddScopedTransformer("foo", UseFieldAsIndex("id")),
		AddScopedTransformer("bar", UseFieldAsIndex("id")),
	)
	flatted, err := f.Flat(toTest)
	assert.Equal(t, 30, flatted["bar.1.data"])
	assert.Equal(t, 40, flatted["bar.2.data"])
	assert.Nil(t, err)
}
