package flatr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScalar(t *testing.T) {
	toTest := 2
	f := New()
	flatted, err := f.Flat(toTest)
	assert.Equal(t, 2, flatted[""])
	assert.Nil(t, err)
}

func TestStruct(t *testing.T) {
	toTest := &struct{ foo int }{foo: 1}
	f := New()
	flatted, err := f.Flat(toTest)
	assert.Equal(t, toTest, flatted[""])
	assert.Nil(t, err)
}

func TestNil(t *testing.T) {
	toTest := map[string]any{
		"foo":  "bar",
		"nest": nil,
	}
	f := New()
	flatted, err := f.Flat(toTest)
	_, ok := flatted["nest"]
	assert.Equal(t, true, ok)
	assert.Nil(t, err)
}

func TestMap(t *testing.T) {
	toTest := map[string]any{
		"foo": "bar",
		"nest": map[string]any{
			"bar": "baz",
		},
	}
	f := New()
	flatted, err := f.Flat(toTest)
	assert.Equal(t, "bar", flatted["foo"])
	assert.Equal(t, "baz", flatted["nest.bar"])
	assert.Nil(t, err)
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
	flatted, err := f.Flat(toTest)
	assert.Equal(t, "bar", flatted["0.foo"])
	assert.Equal(t, 2, flatted["1"])
	assert.Equal(t, "bar", flatted["2"])
	assert.Nil(t, err)
}

func TestMapWithNestedArray(t *testing.T) {

	toTest := map[string]any{
		"foo": []any{1, 2, 3},
	}

	f := New()
	flatted, err := f.Flat(toTest)
	assert.Equal(t, 1, flatted["foo.0"])
	assert.Equal(t, 2, flatted["foo.1"])
	assert.Equal(t, 3, flatted["foo.2"])
	assert.Nil(t, err)
}

func TestPrefixOption(t *testing.T) {
	toTest := map[string]any{
		"foo": "baz",
		"bar": []any{"foobar", "foobaz"},
	}
	f := New(WithPrefix("pre"))
	flatted, err := f.Flat(toTest)
	assert.Equal(t, "baz", flatted["pre.foo"])
	assert.Equal(t, "foobar", flatted["pre.bar.0"])
	assert.Equal(t, "foobaz", flatted["pre.bar.1"])
	assert.Nil(t, err)
}

func TestSeparatorOption(t *testing.T) {
	toTest := map[string]any{
		"foo": map[string]any{
			"bar": "baz",
		},
		"foobar": []any{"foobaz", "foobuz"},
	}
	f := New(WithSeparator("-"))
	flatted, err := f.Flat(toTest)
	assert.Equal(t, "baz", flatted["foo-bar"])
	assert.Equal(t, "foobaz", flatted["foobar-0"])
	assert.Equal(t, "foobuz", flatted["foobar-1"])
	assert.Nil(t, err)
}

func TestTransformer(t *testing.T) {
	toTest := map[string]any{
		"foo": map[string]any{
			"bar": "baz",
		},
		"faa": map[string]any{
			"bar": "baz",
		},
	}
	trasformer := func(e Entry) (Entry, error) {
		if s, ok := e.V.(string); ok {
			e.V = s + "_transformed"
		}
		return e, nil
	}
	f := New(AddTransformer(trasformer))
	flatted, err := f.Flat(toTest)
	assert.Equal(t, "baz_transformed", flatted["foo.bar"])
	assert.Equal(t, "baz_transformed", flatted["faa.bar"])
	assert.Nil(t, err)
}

func TestScopedTransformer(t *testing.T) {
	toTest := map[string]any{
		"foo": map[string]any{
			"bar": "baz",
		},
		"faa": map[string]any{
			"bar": "baz",
		},
	}
	trasformer := func(e Entry) (Entry, error) {
		e.V = e.V.(string) + "_transformed"
		return e, nil
	}
	f := New(AddScopedTransformer("foo.bar", trasformer))
	flatted, err := f.Flat(toTest)
	assert.Equal(t, "baz_transformed", flatted["foo.bar"])
	assert.Equal(t, "baz", flatted["faa.bar"])
	assert.Nil(t, err)
}

func TestTransformerError(t *testing.T) {
	toTest := map[string]any{
		"foo": "bar",
	}
	trasformer := func(e Entry) (Entry, error) {
		e.V = e.V.(string) + "_transformed"
		return e, fmt.Errorf("error")
	}
	f := New(AddScopedTransformer("foo", trasformer))
	flatted, err := f.Flat(toTest)
	assert.Nil(t, flatted["foo"])
	assert.Equal(t, "error", err.Error())
}

func TestCustomTranformer(t *testing.T) {
	toFlat := map[string]any{
		"foo": "bar",
		"nest": map[string]any{
			"bar": "baz",
		},
		"list": []any{
			map[string]any{
				"id":   1,
				"data": "data1",
			},
			map[string]any{
				"id":   2,
				"data": "data2",
			},
		},
	}
	trasformer := func(e Entry) (Entry, error) {
		if s, ok := e.V.(string); ok {
			e.V = s + "_transformed"
		}
		return e, nil
	}
	f := New(AddTransformer(trasformer))
	flatted, _ := f.Flat(toFlat)

	fmt.Printf("%v", flatted)
}
