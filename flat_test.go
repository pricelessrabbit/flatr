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

func TestMapWithNestedArray(t *testing.T) {

	toTest := map[string]any{
		"foo": []any{1, 2, 3},
	}

	f := New()
	flatted := f.Flat(toTest)
	assert.Equal(t, 1, flatted["foo.0"])
	assert.Equal(t, 2, flatted["foo.1"])
	assert.Equal(t, 3, flatted["foo.2"])
}

func TestPrefixOption(t *testing.T) {
	toTest := map[string]any{
		"foo": "bar",
	}
	f := New(WithPrefix("pre"))
	flatted := f.Flat(toTest)
	assert.Equal(t, "bar", flatted["pre.foo"])
}

func TestSeparatorOption(t *testing.T) {
	toTest := map[string]any{
		"foo": map[string]any{
			"bar": "baz",
		},
	}
	f := New(WithSeparator("-"))
	flatted := f.Flat(toTest)
	assert.Equal(t, "baz", flatted["foo-bar"])
}

func TestTransformer(t *testing.T) {
	toTest := map[string]any{
		"foo": map[string]any{
			"bar": "baz",
		},
	}
	trasformer := func(v any) any {
		return v.(string) + "_transformed"
	}

	f := New(AddTransformer("foo.bar", trasformer))
	flatted := f.Flat(toTest)
	assert.Equal(t, "baz_transformed", flatted["foo.bar"])
}

func TestKeyOverIndex(t *testing.T) {
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
		AddTransformer("foo", UseFieldAsIndex("id")),
		AddTransformer("bar", UseFieldAsIndex("id")),
	)
	flatted := f.Flat(toTest)
	assert.Equal(t, 30, flatted["bar.1.data"])
	assert.Equal(t, 40, flatted["bar.2.data"])
}
