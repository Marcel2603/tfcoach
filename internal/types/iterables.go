package types

import "golang.org/x/sync/syncmap"

type Set[T any] struct {
	m syncmap.Map
}

func (s *Set[T]) Add(elem T) {
	s.m.Store(elem, struct{}{})
}

func (s *Set[T]) Values() []T {
	var result []T
	s.m.Range(func(k, _ interface{}) bool {
		elem, ok := k.(T)
		if !ok {
			return false
		}
		result = append(result, elem)
		return true
	})
	return result
}

func (s *Set[T]) Has(elem T) bool {
	_, ok := s.m.Load(elem)
	return ok
}
