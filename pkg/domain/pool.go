package domain

// Pool represent the domain type of pool.
type Pool[T any] interface {
	Get() T
	Put(d T)
}
