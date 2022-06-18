package mapper

import (
	"context"

	"gorm.io/gorm"
)

type Paginator struct {
	StartId uint `json:"start_id"`
	Limit   int  `json:"limit"`
}

type PageRes[T any] struct {
	Paginator
	More bool `json:"more"`
	Data []T  `json:"data"`
}

type Mapper[T any] interface {
	One(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) (*T, error)
	All(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) ([]T, error)
	Paginate(ctx context.Context, pager Paginator, wrapper func(*gorm.DB) *gorm.DB) (*PageRes[T], error)
	Delete(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) error
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB, updated map[string]interface{}) error
}

type BaseMapper[T any] struct {
	db *gorm.DB
}

func (m BaseMapper[T]) One(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) (*T, error) {
	var res T
	err := wrapper(m.db.WithContext(ctx)).
		First(&res).Error
	return &res, err
}

func (m BaseMapper[T]) All(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) ([]T, error) {
	res := make([]T, 0)
	err := wrapper(m.db.WithContext(ctx)).
		Find(&res).Error
	return res, err
}

func (m BaseMapper[T]) Paginate(ctx context.Context, pager Paginator, wrapper func(*gorm.DB) *gorm.DB) (*PageRes[T], error) {
	res := make([]T, 0, pager.Limit+1)
	err := wrapper(m.db.WithContext(ctx)).
		Where("id>?", pager.StartId).
		Limit(pager.Limit + 1).
		Find(&res).Error
	if err != nil {
		return nil, err
	}
	more := len(res) > pager.Limit
	if more {
		res = res[:pager.Limit]
	}

	return &PageRes[T]{
		Paginator: pager,
		More:      more,
		Data:      res,
	}, nil
}

func (m BaseMapper[T]) Delete(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) error {
	var t T
	return wrapper(m.db.WithContext(ctx)).Delete(&t).Error
}

func (m BaseMapper[T]) Create(ctx context.Context, entity *T) error {
	return m.db.WithContext(ctx).Create(entity).Error
}

func (m BaseMapper[T]) Update(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB, updated map[string]interface{}) error {
	return wrapper(m.db.WithContext(ctx)).Updates(updated).Error
}
