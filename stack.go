package flat

type stack[T any] struct {
	entries []T
}

func (s *stack[T]) push(e T) {
	s.entries = append(s.entries, e)
}

func (s *stack[T]) pop() T {
	r := s.entries[len(s.entries)-1]
	s.entries = s.entries[0 : len(s.entries)-1]
	return r
}

func (s *stack[T]) empty() bool {
	return len(s.entries) == 0
}
