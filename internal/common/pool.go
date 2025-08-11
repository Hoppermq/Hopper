package common

import "sync"

type Pool[T any] struct {
	sync.Pool
}

func NewPool[T any](f func() T) *Pool[T] {
	return &Pool[T]{
		Pool: sync.Pool{
			New: func() any { return f() },
		},
	}
}

func (p *Pool[T]) Get() T {
	return p.Pool.Get().(T)
}

func (p *Pool[T]) Put(d T) {
	p.Pool.Put(d)
}
