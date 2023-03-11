package flatr

import "strconv"

type Flatter struct {
	stack             *stack[entry]
	prefix            string
	separator         string
	scopedTrasformers map[string]Trasformer
}

func New(options ...Option) *Flatter {
	f := &Flatter{
		stack:             &stack[entry]{},
		prefix:            "",
		separator:         ".",
		scopedTrasformers: make(map[string]Trasformer),
	}
	for _, opt := range options {
		opt(f)
	}
	return f
}

type entry struct {
	k    string
	v    any
	h    int
	stop bool
}

func AddScopedTransformer(key string, fn Trasformer) Option {
	return func(f *Flatter) {
		f.scopedTrasformers[key] = fn
	}
}

func (f *Flatter) Flat(toFlat any) (map[string]any, error) {
	var err error
	flatted := make(map[string]any)

	f.stack.push(entry{k: f.prefix, v: toFlat, h: 0, stop: false})

	for !f.stack.empty() {
		e := f.stack.pop()
		fn, ok := f.scopedTrasformers[e.k]
		if ok {
			e, err = fn(e)
		}
		if err != nil {
			return flatted, err
		}
		f.flatmapNode(e, flatted)
	}
	return flatted, nil
}

func (f *Flatter) flatmapNode(e entry, flatted map[string]any) {
	if e.stop == true {
		return
	}
	switch e.v.(type) {
	case map[string]any:
		for k, m := range e.v.(map[string]any) {
			nodeKey := joinKey(e.k, k, f.separator)
			f.stack.push(entry{k: nodeKey, v: m})
		}
	case []any:
		for i, v := range e.v.([]any) {
			nodeKey := joinKey(e.k, strconv.Itoa(i), f.separator)
			f.stack.push(entry{k: nodeKey, v: v})
		}
	default:
		flatted[e.k] = e.v
	}
}

func joinKey(parent, k, separator string) string {
	if parent == "" {
		return k
	}
	return parent + separator + k
}
