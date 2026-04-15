package types

import (
	"sync/atomic"

	"golang.org/x/sync/syncmap"
)

type Set[T any] struct {
	m     syncmap.Map
	count atomic.Int32
}

func (s *Set[T]) Add(elem T) {
	_, loaded := s.m.LoadOrStore(elem, struct{}{})
	if !loaded {
		s.count.Add(1)
	}
}

func (s *Set[T]) Values() []T {
	result := make([]T, 0, s.Len())
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

func (s *Set[T]) Len() int {
	return int(s.count.Load())
}

func (s *Set[T]) Has(elem T) bool {
	_, ok := s.m.Load(elem)
	return ok
}

func (s *Set[T]) Delete(elem T) {
	_, present := s.m.LoadAndDelete(elem)
	if present {
		s.count.Add(-1)
	}
}
