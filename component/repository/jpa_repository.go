package repository

type JpaRepository[T, ID any] interface {
	CrudRepository[T, ID]
	PageAndSortRepository[T, ID]
}
