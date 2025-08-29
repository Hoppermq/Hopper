package common

import "sync"

// Pool represent the pooling system.
type Pool[T any] struct {
	sync.Pool
}

// NewPool create a new pool.
func NewPool[T any](f func() T) *Pool[T] {
	return &Pool[T]{
		Pool: sync.Pool{
			New: func() any { return f() },
		},
	}
}

// Get return the item set in the pool.
func (p *Pool[T]) Get() T {
	return p.Pool.Get().(T)
}

// Put set the item in the pool.
func (p *Pool[T]) Put(d T) {
	p.Pool.Put(d)
}
