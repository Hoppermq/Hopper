package domain

type Pool[T any] interface {
	Get() T
	Put(d T)
}
