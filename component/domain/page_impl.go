package domain

type PageImpl[T any] struct {
	total    int
	pageable Pageable
	content  []T
}

func (p PageImpl[T]) GetTotalPages() int {
	if p.GetSize() == 0 {
		return 1
	}
	if p.total == 0 {
		return 0
	}
	return (p.total-1)/p.GetSize() + 1
}

func (p PageImpl[T]) GetTotalElements() int {
	return p.total
}

func (p PageImpl[T]) GetSize() int {
	if p.pageable.IsPaged() {
		return p.pageable.GetPageSize()
	}
	return len(p.content)
}
