package flat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScalar(t *testing.T) {
	toTest := 2
	f := New()
	flatted := f.Flat(toTest)
	assert.Equal(t, 2, flatted[""])
}

func TestStruct(t *testing.T) {
	toTest := &struct{ foo int }{foo: 1}
	f := New()
	flatted := f.Flat(toTest)
	assert.Equal(t, toTest, flatted[""])
}

func TestMap(t *testing.T) {
	toTest := map[string]any{
		"a": "a",
		"nest": map[string]any{
			"b": "b",
		},
	}
	f := New()
	flatted := f.Flat(toTest)
	assert.Equal(t, "a", flatted["a"])
	assert.Equal(t, "b", flatted["nest.b"])
}

func TestArray(t *testing.T) {

	toTest := []any{
		map[string]any{
			"foo": "bar",
		},
		2,
		"bar",
	}

	f := New()
	flatted := f.Flat(toTest)
	assert.Equal(t, "bar", flatted["0.foo"])
	assert.Equal(t, 2, flatted["1"])
	assert.Equal(t, "bar", flatted["2"])
}

func TestMapWithInternalArray(t *testing.T) {

	toTest := map[string]any{
		"foo": []any{1, 2, 3},
	}

	f := New()
	flatted := f.Flat(toTest)
	assert.Equal(t, 1, flatted["foo.0"])
	assert.Equal(t, 2, flatted["foo.1"])
	assert.Equal(t, 3, flatted["foo.2"])
}
