package rx

import "context"

type Node[T any] struct {
	wait  chan struct{}
	next  *Node[T]
	set   bool
	value T
}

func newNode[T any](value *T) *Node[T] {
	if value == nil {
		value = new(T)
	}

	return &Node[T]{
		wait:  make(chan struct{}),
		set:   value != nil,
		value: *value,
	}
}

func (n *Node[T]) iter(ctx context.Context) <-chan T {
	c := make(chan T)
	if ctx == nil {
		ctx = context.Background()
	}
	go func() {
		defer close(c)
		for n := n; ctx.Err() == nil; n = n.next {
			if n.set {
				select {
				case <-ctx.Done():
					return
				case c <- n.value:
				}
			}
			select {
			case <-ctx.Done():
				return
			case <-n.wait:
			}
		}
	}()
	return c
}
