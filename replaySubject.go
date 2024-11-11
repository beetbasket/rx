package rx

import (
	"context"
	"sync"
)

type ReplaySubject[T any] struct {
	lock sync.RWMutex
	head *Node[T]
	tail *Node[T]
}

func (rs *ReplaySubject[T]) Reset() {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.reset()
}

func (rs *ReplaySubject[T]) Append(value ...T) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	for _, v := range value {
		rs.append(newNode[T](&v), false)
	}
}

func (rs *ReplaySubject[T]) Iter(ctx context.Context) <-chan T {
	rs.lock.RLock()
	defer rs.lock.RUnlock()
	if rs.head == nil {
		rs.reset()
	}
	return rs.head.iter(ctx)
}

func (rs *ReplaySubject[T]) reset() {
	rs.append(newNode[T](nil), true)
}

func (rs *ReplaySubject[T]) append(next *Node[T], replaceHead bool) {
	if rs.head == nil || replaceHead {
		rs.head = next
	}
	if rs.tail != nil {
		rs.tail.next = next
		defer close(rs.tail.wait)
	}
	rs.tail = next
}
