<p align="center">
  <img src="https://user-images.githubusercontent.com/22039194/224507070-5e534128-a350-421a-8cb8-3bef2ecce729.png" /><br>
  <i>The Flat Rabbit,  Barour Oskarsson</i>
</p>

# FLATR

Golang util to easily flat a nested map into a 1-level deep one.

## Features

- Prefix selection
- Separator selection
- Transformers to preprocess each node for advanced use cases

## Get it

`go get github.com/pricelessrabbit/flatr`

## Use it

### Basic usage

```go
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
f := flatr.New()
flatted, err := f.Flat(toFlat)


/* result
map[
  foo: bar 
  list.0.data: data1 
  list.0.id: 1 
  list.1.data: data2 
  list.1.id: 2 
  nest.bar: baz
]
*/
```

### Options

Add options to the constructor to use them: `flatr.New(Option1(),Option2()...)`

#### Prefix option

Adds a prefix to all the flatted keys

```go
f := New(flatr.WithPrefix("namespace"))
flatted, _ := f.Flat(toFlat)


/* result
map[
  namespace.foo: bar 
  namespace.list.0.data: data1 
  namespace.list.0.id: 1 
  namespace.list.1.data: data2 
  namespace.list.1.id: 2 
  namespace.nest.bar: baz
]
*/
```


#### Separator option

Choose path separator (default `.`)


```go
f := New(flatr.WithSeparator("_"))
flatted, _ := f.Flat(toFlat)

/* result
map[
  foo: bar
  list_0_data: data1
  list_0_id: 1
  list_1_data: data2
  list_1_id: 2
  nest_bar: baz
]
*/
```

### Transformers

Transformers:
- Have access to the entry related to node being flatted
- Can read and update node key, value, and set a stop flag  to block children processing

#### Default transformers

##### `flatr.Transformers.MaxDeep` 
define a max deep for the flattening process.
after reaching the max deep, the children of the current will be left untouched.

```go
 toFlat := map[string]any{
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
    AddTransformer(MaxDeep(2)),
)
flatted, _ := f.Flat(toFlat)
/*
map[
	foo.0: {"id": "a", "data": 10}
	foo.1: {"id": 1, "data": 30}
	bar.0: {"id": 1, "data": 30}
    bar.1: {"id": 2, "data": 40}
]
*/
	
```

##### `flatr.Transformers.UseFieldAsIndex` 
if the current node is a slice, use the value of the given field (eg `id`) as index instead
of the array index.

```go
toFlat := map[string]any{
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

flatted, _ := f.Flat(toFlat)

/*
map[
foo.a.id: a
foo.a.data: 10
foo.b.id: b 
foo.b.data: 20

bar.1.id: 1
bar.1: data: 30
bar.2.id: 2
bar.2.data: 40
]
*/

```

#### Custom transfomers 
It is possible to implement custom transformers. A trasformer is valid if: 
- takes an Entry
- returns the transformed Entry and an optional error

The entry field that can be transformed are:
- K    the key of the node
- V    the value related to the key 
- H    deep of the node
- Stop stops processing of the children of the current node

##### Example
the transformer in the example adds the suffix "_transformed" to all the string values
that are found in the structure. other value types are left untouched.

```go
customTransformer := func (e Entry) (Entry, error) {
    //check if the current entry value is a string
    // and in that case add the suffix "transformed"
    if s, ok := e.V.(string); ok {
    e.V = s + "_transformed"
    }
    return e, nil
}

f := New(AddTransformer(customTransformer))
flatted, _ := f.Flat(toFlat)

/*
map[
	foo: bar_transformed
	list_0_data: data1_transformed
	list_0_id: 1
	list_1_data: data2_transformed
	list_1_id: 2
	nest_bar: baz_transformed
]
*/

```