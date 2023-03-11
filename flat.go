package flatr

import "strconv"

type Flatter struct {
	stack             *stack[Entry]
	prefix            string
	separator         string
	transformers      []Transformer
	scopedTrasformers map[string]Transformer
}

func New(options ...Option) *Flatter {
	f := &Flatter{
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

type Entry struct {
	K    string
	V    any
	H    int
	Stop bool
}

func (f *Flatter) Flat(toFlat any) (map[string]any, error) {
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

func (f *Flatter) transformEntry(e Entry, err error) (Entry, error) {
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

func (f *Flatter) flatmapNode(e Entry, flatted map[string]any) {
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
