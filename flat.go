package flatr

import "strconv"

type Flatter struct {
	stack       *stack[entry]
	prefix      string
	separator   string
	trasformers map[string]Trasformer
}

func New(options ...Option) *Flatter {
	f := &Flatter{
		stack:       &stack[entry]{},
		prefix:      "",
		separator:   ".",
		trasformers: make(map[string]Trasformer),
	}
	for _, opt := range options {
		opt(f)
	}
	return f
}

type entry struct {
	k string
	v any
}

func (f *Flatter) Flat(toFlat any) map[string]any {

	flatted := make(map[string]any)

	f.stack.push(entry{k: f.prefix, v: toFlat})

	for !f.stack.empty() {
		e := f.stack.pop()
		f.flatmapNode(e.k, e.v, flatted)
	}
	return flatted
}

func (f *Flatter) flatmapNode(rootKey string, toFlat any, flatted map[string]any) {

	fn, ok := f.trasformers[rootKey]
	if ok {
		toFlat = fn(toFlat)
	}
	if toFlat == nil {
		return
	}
	switch toFlat.(type) {
	case map[string]any:
		for k, m := range toFlat.(map[string]any) {
			nodeKey := joinKey(rootKey, k, f.separator)
			f.stack.push(entry{k: nodeKey, v: m})
		}
	case []any:
		for i, v := range toFlat.([]any) {
			nodeKey := joinKey(rootKey, strconv.Itoa(i), f.separator)
			f.stack.push(entry{k: nodeKey, v: v})
		}
	default:
		flatted[rootKey] = toFlat
	}
}

func joinKey(parent, k, separator string) string {
	if parent == "" {
		return k
	}
	return parent + separator + k
}
