package flatr

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

func TestNil(t *testing.T) {
	toTest := map[string]any{
		"foo":  "bar",
		"nest": nil,
	}
	f := New()
	flatted := f.Flat(toTest)
	_, ok := flatted["nest"]
	assert.Equal(t, false, ok)
}

func TestMap(t *testing.T) {
	toTest := map[string]any{
		"foo": "bar",
		"nest": map[string]any{
			"bar": "baz",
		},
	}
	f := New()
	flatted := f.Flat(toTest)
	assert.Equal(t, "bar", flatted["foo"])
	assert.Equal(t, "baz", flatted["nest.bar"])
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
		"foo": "baz",
		"bar": []any{"foobar", "foobaz"},
	}
	f := New(WithPrefix("pre"))
	flatted := f.Flat(toTest)
	assert.Equal(t, "bar", flatted["pre.foo"])
	assert.Equal(t, "foobar", flatted["pre.bar.0"])
	assert.Equal(t, "foobaz", flatted["pre.bar.1"])
}

func TestSeparatorOption(t *testing.T) {
	toTest := map[string]any{
		"foo": map[string]any{
			"bar": "baz",
		},
		"foobar": []any{"foobaz", "foobuz"},
	}
	f := New(WithSeparator("-"))
	flatted := f.Flat(toTest)
	assert.Equal(t, "baz", flatted["foo-bar"])
	assert.Equal(t, "foobaz", flatted["foobar-0"])
	assert.Equal(t, "foobuz", flatted["foobar-1"])
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
		AddTransformer("foo", UseFieldAsIndex("id")),
		AddTransformer("bar", UseFieldAsIndex("id")),
	)
	flatted := f.Flat(toTest)
	assert.Equal(t, 30, flatted["bar.1.data"])
	assert.Equal(t, 40, flatted["bar.2.data"])
}
