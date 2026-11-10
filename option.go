package main

type Option[T any] struct {
	value T
	valid bool
}

func Some[T any](Value T) Option[T] {
	return Option[T]{value: Value, valid: true}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func (o Option[T]) IsSome() bool {
	return o.valid
}

func (o Option[T]) IsNone() bool {
	return !o.valid
}

func (o Option[T]) Unwarp() T {
	if !o.valid {
		panic("called unwarp on a None value")
	}
	return o.value
}

func (o Option[T]) UnwarpOr(def T) T {
	if o.valid {
		return o.value
	}
	return def
}

func (o Option[T]) UnwarpOrElse(fn func() T) T {
	if o.valid {
		return o.value
	}
	return fn()
}
