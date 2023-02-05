package flatthis

type Flatter struct {
	stack stack
}

func New() *Flatter {
	return &Flatter{}
}

type entry struct {
	k string
	v map[string]any
}

type stack struct {
	entries []entry
}

func (s *stack) push(e entry) {
	s.entries = append(s.entries, e)
}

func (s *stack) pop() entry {
	r := s.entries[len(s.entries)-1]
	s.entries = s.entries[0 : len(s.entries)-1]
	return r
}

func (s *stack) empty() bool {
	return len(s.entries) == 0
}

func (f *Flatter) FlatMap(toFlat map[string]any) map[string]any {
	s := &stack{}
	flatted := make(map[string]any)

	s.push(entry{k: "", v: toFlat})

	for !s.empty() {
		e := s.pop()
		flatmapNode(e.k, e.v, flatted, s)
	}
	return flatted
}

func flatmapNode(rootKey string, toFlat map[string]any, flatted map[string]any, stack *stack) map[string]any {
	for k, m := range toFlat {
		nodeKey := joinKey(rootKey, k, ".")
		switch m.(type) {
		case map[string]any:
			stack.push(entry{k: nodeKey, v: m.(map[string]any)})
		default:
			flatted[nodeKey] = m
		}
	}
	return flatted
}

func joinKey(root, k, separator string) string {
	if root == "" {
		return k
	}
	return root + separator + k
}
