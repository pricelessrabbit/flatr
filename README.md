<p align="center">
  <img src="https://user-images.githubusercontent.com/22039194/224507070-5e534128-a350-421a-8cb8-3bef2ecce729.png" /><br>
  <i>The Rlat Rabbit,  Barour Oskarsson</i>
</p>

# FLATR

Golang util to easily flat a nested map into a 1-level deep one.

## Features

* Prefix selection
* Separator selection
* Transformers to preprocess each node for advanced use cases

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
* has access to the entry related to node being flatted
* can read and update node key, value, and set a stop flag  to block children processing

#### Basic transfomer example

```go
trasformer := func(e Entry) (Entry, error) {
		//check if the current entry value is a string
    // and in that case add the suffix "transformed"
    if s, ok := e.V.(string); ok {
			e.V = s + "_transformed"
		}
		return e, nil
	}
  
f := New(AddTransformer(trasformer))
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
