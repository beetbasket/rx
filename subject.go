package rx

import (
	"context"
	"sync"
)

type Subject[T any] struct {
	lock sync.RWMutex
	init sync.Once
	head *node[T]
}

func (s *Subject[T]) Next(value ...T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, v := range value {
		s.append(&v)
	}
	s.append(nil)
}

func (s *Subject[T]) Subscribe(ctx context.Context) <-chan T {
	s.initFn()
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.head.iter(ctx)
}

func (s *Subject[T]) Complete(value ...T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, v := range value {
		s.append(&v)
	}
	s.appendFinal()
}

func (s *Subject[T]) initFn() {
	s.init.Do(func() {
		s.lock.Lock()
		defer s.lock.Unlock()
		if s.head == nil {
			s.head = newNode[T](nil)
		}
	})
}

func (s *Subject[T]) append(v *T) {
	if s.head != nil && s.head.final {
		return
	}

	n := newNode(v)
	if s.head != nil {
		s.head.next = n
		close(s.head.wait)
	}
	s.head = n
}

func (s *Subject[T]) appendFinal() {
	s.append(nil)
	if s.head.final {
		return
	}
	n := newNode[T](nil)
	n.final = true
	s.head.next = n
	close(s.head.wait)
	s.head = n
}
