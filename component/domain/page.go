package domain

type Page[T any] interface {
	GetTotalPages() int
	GetTotalElements() int
}
