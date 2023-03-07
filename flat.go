package flat

import "strconv"

type Flatter struct {
	stack stack[entry]
}

func New() *Flatter {
	return &Flatter{}
}

type entry struct {
	k string
	v any
}

func (f *Flatter) Flat(toFlat any) map[string]any {
	s := &stack[entry]{}
	flatted := make(map[string]any)

	s.push(entry{k: "", v: toFlat})

	for !s.empty() {
		e := s.pop()
		flatmapNode(e.k, e.v, flatted, s)
	}
	return flatted
}

func flatmapNode(rootKey string, toFlat any, flatted map[string]any, stack *stack[entry]) {
	switch toFlat.(type) {
	case map[string]any:
		for k, m := range toFlat.(map[string]any) {
			nodeKey := joinKey(rootKey, k, ".")
			stack.push(entry{k: nodeKey, v: m.(any)})
		}
	case []any:
		for i, v := range toFlat.([]any) {
			nodeKey := joinKey(rootKey, strconv.Itoa(i), ".")
			stack.push(entry{k: nodeKey, v: v.(any)})
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
