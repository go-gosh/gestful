package repository

import "github.com/go-gosh/gestful/component/domain"

type PageAndSortRepository[T, ID any] interface {
	Repository[T, ID]
	FindAllBySort(sort domain.Sort) ([]T, error)
	FindAllByPage(page domain.Pageable) (domain.Page[T], error)
}
