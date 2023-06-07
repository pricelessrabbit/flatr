package flatr

import "strconv"

// Flattener keep memory of options and transformers to flatten maps
type Flattener struct {
	stack             *stack[Entry]
	prefix            string
	separator         string
	transformers      []Transformer
	scopedTrasformers map[string]Transformer
}

// New instantiate new flatter with provided options
// WithPrefix option add a prefix to keys
// WithSeparator option change default separator
// AddTransformer and AddScopedTransformer add custom transfomers that apply
// a custom function to each entry before processing it
func New(options ...Option) *Flattener {
	f := &Flattener{
		stack:             &stack[Entry]{},
		prefix:            "",
		separator:         ".",
		scopedTrasformers: make(map[string]Transformer),
	}
	for _, opt := range options {
		opt(f)
	}
	return f
}

// Entry keep intermediate structure when flatten is ongoing and is passed to transformers
// tracked data include key and value of the node, height in the tree and if process have to
// stop without processing children nodes
type Entry struct {
	K    string
	V    any
	H    int
	Stop bool
}

// Flat flatten a nested data structure (scalar, maps and slices) to a flat map.
// returns a map of 1 level when key is the full path of the field divided by a separator (default to .)
func (f *Flattener) Flat(toFlat any) (map[string]any, error) {
	var err error
	flatted := make(map[string]any)

	f.stack.push(Entry{K: f.prefix, V: toFlat, H: 0, Stop: false})

	for !f.stack.empty() {
		e := f.stack.pop()

		e, err = f.transformEntry(e, err)
		if err != nil {
			return flatted, err
		}
		f.flatmapNode(e, flatted)
	}
	return flatted, nil
}

func (f *Flattener) transformEntry(e Entry, err error) (Entry, error) {
	transformers := f.transformers
	fn, ok := f.scopedTrasformers[e.K]
	if ok {
		transformers = append(transformers, fn)
	}
	for _, fn := range transformers {
		e, err = fn(e)
		if err != nil {
			return e, err
		}
	}
	return e, nil
}

func (f *Flattener) flatmapNode(e Entry, flatted map[string]any) {
	if e.Stop == true {
		flatted[e.K] = e.V
		return
	}
	switch e.V.(type) {
	case map[string]any:
		for k, m := range e.V.(map[string]any) {
			nodeKey := joinKey(e.K, k, f.separator)
			f.stack.push(Entry{K: nodeKey, V: m, H: e.H + 1})
		}
	case []any:
		for i, v := range e.V.([]any) {
			nodeKey := joinKey(e.K, strconv.Itoa(i), f.separator)
			f.stack.push(Entry{K: nodeKey, V: v, H: e.H + 1})
		}
	default:
		flatted[e.K] = e.V
	}
}

func joinKey(parent, k, separator string) string {
	if parent == "" {
		return k
	}
	return parent + separator + k
}
