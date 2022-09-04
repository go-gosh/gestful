package mapper

import (
	"context"

	"gorm.io/gorm"
)

type IMapper[T, U, V any] interface {
	IQueryMapper[T, U, V]
	ICommandMapper[T]
}

type IQueryMapper[T, U, V any] interface {
	One(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) (*T, error)
	OneById(ctx context.Context, id uint) (*T, error)
	All(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) ([]T, error)
	Paginate(ctx context.Context, pager U, wrapper func(*gorm.DB) *gorm.DB) (*V, error)
	Count(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) (int, error)
}

type ICommandMapper[T any] interface {
	Delete(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) error
	DeleteById(ctx context.Context, id uint) error
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB, updated map[string]interface{}) error
	UpdateById(ctx context.Context, id uint, updated map[string]interface{}) error
}
