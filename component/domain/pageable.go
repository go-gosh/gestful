package domain

type Pageable interface {
	IsPaged() bool
	GetPageNumber() int
	GetPageSize() int
	GetOffset() int
	GetSort() Sort
}
