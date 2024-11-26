package rx

import "context"

type node[T any] struct {
	wait  chan struct{}
	next  *node[T]
	set   bool
	value T
	final bool
}

func newNode[T any](value *T) *node[T] {
	v := value
	if v == nil {
		v = new(T)
	}

	return &node[T]{
		wait:  make(chan struct{}),
		set:   value != nil,
		value: *v,
	}
}

func (n *node[T]) iter(ctx context.Context) <-chan T {
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
			} else if n.final {
				return
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
