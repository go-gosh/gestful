package repository

type CrudRepository[T, ID any] interface {
	Repository[T, ID]
	Save(entity T) (T, error)
	SaveAll(entity ...T) ([]T, error)
	FindById(id ID) (*T, error)
	ExistsById(id ID) (bool, error)
	FindAll() ([]T, error)
	FindAllById(id ...[]ID) ([]T, error)
	Count() (int, error)
	DeleteById(id ID) error
	Delete(entity T) error
	DeleteAllById(id ...[]ID) error
	DeleteAll(entity ...T) error
}
